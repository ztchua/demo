const API = '/expenses';
let editMode = false;
let initialLoad = true;
let pendingDeleteId = null;

const form = document.getElementById('expenseForm');
const formTitle = document.getElementById('formTitle');
const cancelBtn = document.getElementById('cancelBtn');
const expenseList = document.getElementById('expenseList');

// Load expenses on page load
fetchExpenses();

form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const dateInput = document.getElementById('date').value;
    let dateValue;

    if (dateInput) {
        // datetime-local returns format "YYYY-MM-DDTHH:mm" - need to convert to ISO
        const d = new Date(dateInput);
        dateValue = d.toISOString();
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

    try {
        if (editMode && id) {
            await updateExpense(id, data);
            showToast('Expense updated successfully');
        } else {
            await createExpense(data);
            showToast('Expense added successfully');
        }

        resetForm();
        await fetchExpenses();
    } catch (error) {
        console.error('Error saving expense:', error);
        showToast('Failed to save expense', true);
    }
});

cancelBtn.addEventListener('click', resetForm);

async function fetchExpenses() {
    const res = await fetch(API);
    const expenses = await res.json();
    renderExpenses(expenses);
}

function renderExpenses(expenses) {
    if (expenses.length === 0) {
        expenseList.innerHTML = `
            <tr>
                <td colspan="5">
                    <div class="empty-state">
                        No expenses yet. Add your first expense above!
                    </div>
                </td>
            </tr>
        `;
        initialLoad = false;
        return;
    }

    expenseList.innerHTML = expenses.map((e, index) => `
        <tr style="animation: rowSlideIn 0.4s ease-out forwards ${initialLoad ? (0.3 + (index * 0.05)) : 0}s">
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

    initialLoad = false;
}

async function createExpense(data) {
    const res = await fetch(API, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    if (!res.ok) {
        throw new Error(`Failed to create expense: ${res.status}`);
    }
    return await res.json();
}

async function updateExpense(id, data) {
    const res = await fetch(`${API}/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    });
    if (!res.ok) {
        throw new Error(`Failed to update expense: ${res.status}`);
    }
}

async function deleteExpense(id) {
    const row = findRowByExpenseId(id);

    if (row) {
        row.classList.add('row-deleting');
        await new Promise(resolve => setTimeout(resolve, 400));
    }

    await fetch(`${API}/${id}`, { method: 'DELETE' });
    fetchExpenses();
    showToast('Expense deleted');
}

function findRowByExpenseId(id) {
    const rows = expenseList.querySelectorAll('tr');
    for (const row of rows) {
        const editBtn = row.querySelector(`button[onclick="editExpense(${id})"]`);
        if (editBtn) return row;
    }
    return null;
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

            // Add edit mode styling
            document.querySelector('.card-form').classList.add('edit-mode');
            cancelBtn.classList.add('visible');
        });
}

function confirmDelete(id) {
    pendingDeleteId = id;
    const modal = document.getElementById('deleteModal');
    modal.classList.add('visible');
}

// Modal handlers
document.getElementById('modalCancel').addEventListener('click', () => {
    const modal = document.getElementById('deleteModal');
    modal.classList.remove('visible');
    pendingDeleteId = null;
});

document.getElementById('modalConfirm').addEventListener('click', async () => {
    if (pendingDeleteId) {
        await deleteExpense(pendingDeleteId);
        pendingDeleteId = null;
    }
    document.getElementById('deleteModal').classList.remove('visible');
});

// Close modal on backdrop click
document.getElementById('deleteModal').addEventListener('click', (e) => {
    if (e.target.id === 'deleteModal') {
        e.target.classList.remove('visible');
        pendingDeleteId = null;
    }
});

function resetForm() {
    form.reset();
    document.getElementById('expenseId').value = '';
    editMode = false;
    formTitle.textContent = 'Add Expense';

    // Remove edit mode styling
    document.querySelector('.card-form').classList.remove('edit-mode');
    cancelBtn.classList.remove('visible');
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

// Toast notification function
function showToast(message, isError = false) {
    const toast = document.createElement('div');
    toast.className = `toast ${isError ? 'error' : ''}`;
    toast.textContent = message;
    document.body.appendChild(toast);

    // Trigger animation
    requestAnimationFrame(() => {
        toast.classList.add('visible');
    });

    // Auto-remove
    setTimeout(() => {
        toast.classList.remove('visible');
        setTimeout(() => toast.remove(), 400);
    }, 3000);
}
