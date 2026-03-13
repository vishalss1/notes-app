const API_URL = '/notes';

// DOM Elements
const noteForm = document.getElementById('note-form');
const noteIdInput = document.getElementById('note-id');
const noteTitleInput = document.getElementById('note-title');
const noteContentInput = document.getElementById('note-content');
const submitBtn = document.getElementById('submit-btn');
const cancelEditBtn = document.getElementById('cancel-edit-btn');
const notesContainer = document.getElementById('notes-container');
const noteCardTemplate = document.getElementById('note-card-template');

// State
let isEditing = false;

// Event Listeners
document.addEventListener('DOMContentLoaded', fetchNotes);
noteForm.addEventListener('submit', handleFormSubmit);
cancelEditBtn.addEventListener('click', resetForm);

// Fetch all notes
async function fetchNotes() {
    try {
        const response = await fetch(API_URL);
        if (!response.ok) throw new Error('Failed to fetch notes');
        const notes = await response.json();
        renderNotes(notes);
    } catch (error) {
        console.error('Error:', error);
        notesContainer.innerHTML = `<div class="empty-state"><i class="fa-solid fa-triangle-exclamation"></i> Error loading notes. Please try again.</div>`;
    }
}

// Render notes to DOM
function renderNotes(notes) {
    notesContainer.innerHTML = '';
    
    if (!notes || notes.length === 0) {
        notesContainer.innerHTML = `<div class="empty-state"><i class="fa-regular fa-folder-open"></i> No notes found. Create your first note above!</div>`;
        return;
    }

    // Sort notes descending by ID (assuming newer notes have higher IDs)
    notes.sort((a, b) => b.id - a.id);

    notes.forEach(note => {
        const clone = noteCardTemplate.content.cloneNode(true);
        const card = clone.querySelector('.note-card');
        
        card.dataset.id = note.id;
        clone.querySelector('.note-title').textContent = note.title;
        clone.querySelector('.note-content').textContent = note.content;
        
        // Add event listeners to buttons
        const editBtn = clone.querySelector('.edit-btn');
        editBtn.addEventListener('click', () => initiateEdit(note));
        
        const deleteBtn = clone.querySelector('.delete-btn');
        deleteBtn.addEventListener('click', () => deleteNote(note.id));

        notesContainer.appendChild(clone);
    });
}

// Handle form submission (Create or Update)
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
            noteData.id = parseInt(id, 10);
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

// Create a new note
async function createNote(noteData) {
    const response = await fetch(API_URL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(noteData)
    });
    
    if (!response.ok) throw new Error('Failed to create note');
}

// Update an existing note
async function updateNote(id, noteData) {
    const response = await fetch(`${API_URL}/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(noteData)
    });
    
    if (!response.ok) throw new Error('Failed to update note');
}

// Delete a note
async function deleteNote(id) {
    try {
        const response = await fetch(`${API_URL}/${id}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) throw new Error('Failed to delete note');
        
        fetchNotes();
    } catch (error) {
        console.error('Error deleting note:', error);
        alert('Failed to delete note. Please try again.');
    }
}

// Initiate editing mode
function initiateEdit(note) {
    isEditing = true;
    noteIdInput.value = note.id;
    noteTitleInput.value = note.title;
    noteContentInput.value = note.content;
    
    submitBtn.innerHTML = '<i class="fa-solid fa-check"></i> Update Note';
    cancelEditBtn.classList.remove('hidden');
    
    // Scroll to form smoothly
    noteForm.scrollIntoView({ behavior: 'smooth', block: 'center' });
    noteTitleInput.focus();
}

// Reset form to default state
function resetForm() {
    isEditing = false;
    noteForm.reset();
    noteIdInput.value = '';
    
    submitBtn.innerHTML = '<i class="fa-solid fa-plus"></i> Add Note';
    cancelEditBtn.classList.add('hidden');
}
