const API = '/expenses';
let editMode = false;

const form = document.getElementById('expenseForm');
const formTitle = document.getElementById('formTitle');
const cancelBtn = document.getElementById('cancelBtn');
const expenseList = document.getElementById('expenseList');

// Load expenses on page load
fetchExpenses();

form.addEventListener('submit', async (e) => {
    e.preventDefault();

    // Handle date formatting for datetime-local input
    const dateInput = document.getElementById('date').value;
    let dateValue;
    if (dateInput) {
        // datetime-local returns "YYYY-MM-DDTHH:mm" - convert to ISO format
        dateValue = new Date(dateInput).toISOString();
    } else {
        dateValue = new Date().toISOString();
    }

    const data = {
        description: document.getElementById('description').value,
        amount: parseFloat(document.getElementById('amount').value),
        category: document.getElementById('category').value,
        date: dateValue
    };

    const id = document.getElementById('expenseId').value;

    if (editMode && id) {
        await updateExpense(id, data);
    } else {
        await createExpense(data);
    }

    resetForm();
    fetchExpenses();
});

cancelBtn.addEventListener('click', resetForm);

async function fetchExpenses() {
    const res = await fetch(API);
    const expenses = await res.json();
    renderExpenses(expenses);
}

function renderExpenses(expenses) {
    expenseList.innerHTML = expenses.map(e => `
        <tr>
            <td>${escapeHtml(e.description)}</td>
            <td>$${e.amount.toFixed(2)}</td>
            <td>${escapeHtml(e.category || '-')}</td>
            <td>${formatDate(e.date)}</td>
            <td class="actions">
                <button class="btn-secondary btn-small" onclick="editExpense(${e.id})">Edit</button>
                <button class="btn-danger btn-small" onclick="confirmDelete(${e.id})">Delete</button>
            </td>
        </tr>
    `).join('');
}

async function createExpense(data) {
    await fetch(API, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
}

async function updateExpense(id, data) {
    await fetch(`${API}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
}

async function deleteExpense(id) {
    await fetch(`${API}/${id}`, { method: 'DELETE' });
    fetchExpenses();
}

function editExpense(id) {
    fetch(`${API}/${id}`)
        .then(res => res.json())
        .then(expense => {
            document.getElementById('expenseId').value = expense.id;
            document.getElementById('description').value = expense.description;
            document.getElementById('amount').value = expense.amount;
            document.getElementById('category').value = expense.category || '';
            document.getElementById('date').value = expense.date?.slice(0, 16) || '';

            editMode = true;
            formTitle.textContent = 'Edit Expense';
            cancelBtn.style.display = 'inline-block';
        });
}

function confirmDelete(id) {
    if (confirm('Delete this expense?')) {
        deleteExpense(id);
    }
}

function resetForm() {
    form.reset();
    document.getElementById('expenseId').value = '';
    editMode = false;
    formTitle.textContent = 'Add Expense';
    cancelBtn.style.display = 'none';
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
}
