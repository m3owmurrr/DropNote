async function sendNote() {
    const text = document.getElementById('noteText').value;
    const isPublic = document.getElementById('isPublic').checked;

    const response = await fetch('http://localhost:8080/api/notes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ text: text, public: isPublic })
    });

    if (response.ok) {
        const result = await response.json();
        document.getElementById('noteId').value = result.note_id;
    } else {
        alert('Ошибка при отправке заметки');
    }
}

async function findNote() {
    const noteId = document.getElementById('noteId').value;
    if (!noteId) {
        alert('Введите ID заметки');
        return;
    }

    const response = await fetch(`http://localhost:8080/api/notes/${noteId}`);
    if (response.ok) {
        const result = await response.json();
        document.getElementById('noteText').value = result.text;
    } else {
        alert('Ошибка при получении заметки');
    }
}
