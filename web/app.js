const API = "/notes";

async function loadNotes() {
  const res = await fetch("/notes");

  if (!res.ok) {
    console.error("Failed to load notes");
    return;
  }

  const notes = await res.json();

  if (!Array.isArray(notes)) {
    console.error("Expected array, got:", notes);
    return;
  }

  const list = document.getElementById("notesList");
  list.innerHTML = "";

  notes.forEach(note => {
    const li = document.createElement("li");
    li.textContent = `${note.id}: ${note.title} - ${note.content}`;
    list.appendChild(li);
  });
}

async function createNote() {
  const title = document.getElementById("title").value;
  const content = document.getElementById("content").value;

  if (!title) {
    alert("Title is required");
    return;
  }

  await fetch(API, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ title, content })
  });

  document.getElementById("title").value = "";
  document.getElementById("content").value = "";

  loadNotes();
}

async function deleteNote(id) {
  await fetch(`${API}/${id}`, {
    method: "DELETE"
  });

  loadNotes();
}

async function editNote(id) {
  const newTitle = prompt("New title:");
  const newContent = prompt("New content:");

  if (!newTitle) return;

  await fetch(`${API}/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      id,
      title: newTitle,
      content: newContent
    })
  });

  loadNotes();
}

window.onload = loadNotes;