const API_URL = '/notes';
const AUTH_URL = '/auth';

// DOM Elements
const authSection = document.getElementById('auth-section');
const appSection = document.getElementById('app-section');
const authForm = document.getElementById('auth-form');
const authEmailInput = document.getElementById('auth-email');
const authPasswordInput = document.getElementById('auth-password');
const authSubmitBtn = document.getElementById('auth-submit-btn');
const authToggleBtn = document.getElementById('auth-toggle-btn');
const authToggleText = document.getElementById('auth-toggle-text');
const authSubtitle = document.getElementById('auth-subtitle');
const authError = document.getElementById('auth-error');
const logoutBtn = document.getElementById('logout-btn');

const noteForm = document.getElementById('note-form');
const noteIdInput = document.getElementById('note-id');
const noteTitleInput = document.getElementById('note-title');
const noteContentInput = document.getElementById('note-content');
const submitBtn = document.getElementById('submit-btn');
const cancelEditBtn = document.getElementById('cancel-edit-btn');
const notesContainer = document.getElementById('notes-container');
const noteCardTemplate = document.getElementById('note-card-template');

// State
let accessToken = null; // In-memory only — never in localStorage
let isEditing = false;
let isLoginMode = true;
let isRefreshing = false;
let refreshPromise = null;

// Event Listeners
document.addEventListener('DOMContentLoaded', init);
authForm.addEventListener('submit', handleAuth);
authToggleBtn.addEventListener('click', toggleAuthMode);
logoutBtn.addEventListener('click', handleLogout);
noteForm.addEventListener('submit', handleFormSubmit);
cancelEditBtn.addEventListener('click', resetForm);

// Initialization
async function init() {
    const token = await tryRefreshToken();
    if (token) {
        showApp();
    } else {
        showAuth();
    }
}

// Auth API wrapper with auto-refresh
async function apiCall(url, options = {}) {
    if (!options.headers) options.headers = {};
    if (accessToken) {
        options.headers['Authorization'] = `Bearer ${accessToken}`;
    }

    let response = await fetch(url, options);

    if (response.status === 401 && accessToken) {
        const refreshed = await tryRefreshToken();
        if (refreshed) {
            options.headers['Authorization'] = `Bearer ${accessToken}`;
            response = await fetch(url, options);
        } else {
            handleSessionExpired();
            throw new Error('Session expired');
        }
    }

    return response;
}

async function tryRefreshToken() {
    if (isRefreshing) {
        return refreshPromise;
    }

    isRefreshing = true;
    refreshPromise = (async () => {
        try {
            const response = await fetch(`${AUTH_URL}/refresh`, {
                method: 'POST',
                credentials: 'same-origin',
            });

            if (!response.ok) {
                accessToken = null;
                return false;
            }

            const data = await response.json();
            accessToken = data.access_token;
            return true;
        } catch {
            accessToken = null;
            return false;
        } finally {
            isRefreshing = false;
            refreshPromise = null;
        }
    })();

    return refreshPromise;
}

// Auth handlers
async function handleAuth(e) {
    e.preventDefault();
    hideAuthError();

    const email = authEmailInput.value.trim();
    const password = authPasswordInput.value;

    if (!email || !password) return;

    const endpoint = isLoginMode ? 'login' : 'signup';

    try {
        authSubmitBtn.disabled = true;
        authSubmitBtn.innerHTML = '<i class="fa-solid fa-circle-notch fa-spin"></i> Please wait...';

        const response = await fetch(`${AUTH_URL}/${endpoint}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'same-origin',
            body: JSON.stringify({ email, password }),
        });

        if (!response.ok) {
            const errorText = await response.text();
            let errorMsg;
            try {
                errorMsg = JSON.parse(errorText).error;
            } catch {
                errorMsg = errorText;
            }
            showAuthError(errorMsg || 'Authentication failed');
            return;
        }

        const data = await response.json();
        accessToken = data.access_token;
        showApp();
    } catch (error) {
        showAuthError('Network error. Please try again.');
    } finally {
        authSubmitBtn.disabled = false;
        authSubmitBtn.innerHTML = isLoginMode
            ? '<i class="fa-solid fa-right-to-bracket"></i> Login'
            : '<i class="fa-solid fa-user-plus"></i> Sign Up';
    }
}

async function handleLogout() {
    try {
        await fetch(`${AUTH_URL}/logout`, {
            method: 'POST',
            credentials: 'same-origin',
        });
    } catch {
        // Ignore errors on logout
    }

    accessToken = null;
    showAuth();
}

function handleSessionExpired() {
    accessToken = null;
    showAuth();
    showAuthError('Session expired. Please login again.');
}

// UI toggle
function showApp() {
    authSection.classList.add('hidden');
    appSection.classList.remove('hidden');
    authForm.reset();
    hideAuthError();
    fetchNotes();
}

function showAuth() {
    appSection.classList.add('hidden');
    authSection.classList.remove('hidden');
    isLoginMode = true;
    updateAuthUI();
}

function toggleAuthMode() {
    isLoginMode = !isLoginMode;
    updateAuthUI();
    hideAuthError();
}

function updateAuthUI() {
    if (isLoginMode) {
        authSubtitle.textContent = 'Sign in to your account';
        authSubmitBtn.innerHTML = '<i class="fa-solid fa-right-to-bracket"></i> Login';
        authToggleText.textContent = "Don't have an account?";
        authToggleBtn.textContent = 'Sign up';
    } else {
        authSubtitle.textContent = 'Create a new account';
        authSubmitBtn.innerHTML = '<i class="fa-solid fa-user-plus"></i> Sign Up';
        authToggleText.textContent = 'Already have an account?';
        authToggleBtn.textContent = 'Login';
    }
}

function showAuthError(msg) {
    authError.textContent = msg;
    authError.classList.remove('hidden');
}

function hideAuthError() {
    authError.textContent = '';
    authError.classList.add('hidden');
}

// Notes CRUD
async function fetchNotes() {
    try {
        const response = await apiCall(API_URL);
        if (!response.ok) throw new Error('Failed to fetch notes');
        const notes = await response.json();
        renderNotes(notes);
    } catch (error) {
        console.error('Error:', error);
        notesContainer.innerHTML = `<div class="empty-state"><i class="fa-solid fa-triangle-exclamation"></i> Error loading notes. Please try again.</div>`;
    }
}

function renderNotes(notes) {
    notesContainer.innerHTML = '';

    if (!notes || notes.length === 0) {
        notesContainer.innerHTML = `<div class="empty-state"><i class="fa-regular fa-folder-open"></i> No notes found. Create your first note above!</div>`;
        return;
    }

    notes.forEach(note => {
        const clone = noteCardTemplate.content.cloneNode(true);
        const card = clone.querySelector('.note-card');

        card.dataset.id = note.id;
        clone.querySelector('.note-title').textContent = note.title;
        clone.querySelector('.note-content').textContent = note.content;

        const editBtn = clone.querySelector('.edit-btn');
        editBtn.addEventListener('click', () => initiateEdit(note));

        const deleteBtn = clone.querySelector('.delete-btn');
        deleteBtn.addEventListener('click', () => deleteNote(note.id));

        notesContainer.appendChild(clone);
    });
}

async function handleFormSubmit(e) {
    e.preventDefault();

    const noteData = {
        title: noteTitleInput.value.trim(),
        content: noteContentInput.value.trim()
    };

    if (!noteData.title || !noteData.content) return;

    try {
        if (isEditing) {
            const id = noteIdInput.value;
            await updateNote(id, noteData);
        } else {
            await createNote(noteData);
        }

        resetForm();
        fetchNotes();
    } catch (error) {
        console.error('Error saving note:', error);
        alert('Failed to save note. Please try again.');
    }
}

async function createNote(noteData) {
    const response = await apiCall(API_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(noteData)
    });

    if (!response.ok) throw new Error('Failed to create note');
}

async function updateNote(id, noteData) {
    const response = await apiCall(`${API_URL}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(noteData)
    });

    if (!response.ok) throw new Error('Failed to update note');
}

async function deleteNote(id) {
    try {
        const response = await apiCall(`${API_URL}/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete note');

        fetchNotes();
    } catch (error) {
        console.error('Error deleting note:', error);
        alert('Failed to delete note. Please try again.');
    }
}

function initiateEdit(note) {
    isEditing = true;
    noteIdInput.value = note.id;
    noteTitleInput.value = note.title;
    noteContentInput.value = note.content;

    submitBtn.innerHTML = '<i class="fa-solid fa-check"></i> Update Note';
    cancelEditBtn.classList.remove('hidden');

    noteForm.scrollIntoView({ behavior: 'smooth', block: 'center' });
    noteTitleInput.focus();
}

function resetForm() {
    isEditing = false;
    noteForm.reset();
    noteIdInput.value = '';

    submitBtn.innerHTML = '<i class="fa-solid fa-plus"></i> Add Note';
    cancelEditBtn.classList.add('hidden');
}
