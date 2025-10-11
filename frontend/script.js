const apiUrl = '/itens';

const form = document.getElementById('form-item');
const itemIdInput = document.getElementById('item-id');
const nameInput = document.getElementById('name');
const priceInput = document.getElementById('price');
const submitButton = document.getElementById('btn-submit');
const cancelButton = document.getElementById('btn-cancel');

async function loadItems() {
    try {
        const res = await fetch(apiUrl);
        const items = await res.json();

        const ul = document.getElementById('item-list');
        ul.innerHTML = '';

        if (!items || items.length === 0) {
            ul.innerHTML = '<li>No items registered</li>';
            return;
        }

        items.forEach(item => {
            const li = document.createElement('li');
            li.innerHTML = `
                <span><b>ID:</b> ${item.id} — ${item.name}: $ ${item.price.toFixed(2)}</span>
                <div class="item-actions">
                    <button class="btn-edit" onclick='prepareEdit(${JSON.stringify(item)})'>Edit</button>
                    <button class="btn-delete" onclick="deleteItem(${item.id})">Delete</button>
                </div>
            `;
            ul.appendChild(li);
        });
    } catch (err) {
        console.error('Error loading items:', err);
    }
}

form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const id = itemIdInput.value;
    const item = {
        name: nameInput.value,
        price: parseFloat(priceInput.value)
    };

    if (id) {
        await updateItem(id, item);
    } else {
        await addItem(item);
    }
});

async function addItem(item) {
    try {
        const res = await fetch(apiUrl, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(item)
        });
        if (!res.ok) {
            const msg = await res.text();
            alert('Error adding item: ' + msg);
            return;
        }
        form.reset();
        loadItems();
    } catch (err) {
        console.error('Error adding item:', err);
        alert('Error adding item (network): ' + err.message);
    }
}

function prepareEdit(item) {
    itemIdInput.value = item.id;
    nameInput.value = item.name;
    priceInput.value = item.price;

    submitButton.textContent = 'Update Item';
    cancelButton.style.display = 'inline-block';
    window.scrollTo(0, 0);
}

async function updateItem(id, item) {
    try {
        const res = await fetch(`${apiUrl}/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(item)
        });
        if (!res.ok) {
            const msg = await res.text();
            alert('Error updating item: ' + msg);
            return;
        }
        cancelEdit();
        loadItems();
    } catch (err) {
        console.error('Error updating item:', err);
        alert('Error updating item (network): ' + err.message);
    }
}

async function searchItem() {
    const id = document.getElementById('search').value;
    const ulResult = document.getElementById('search-result');
    ulResult.innerHTML = '';

    if (!id) {
        alert("Enter an ID to search.");
        return;
    }

    try {
        const res = await fetch(`${apiUrl}/${id}`);

        if (!res.ok) {
            ulResult.innerHTML = `<li style="color:red;">Item not found</li>`;
            return;
        }

        const item = await res.json();

        const li = document.createElement('li');
        li.innerHTML = `
            <span><b>ID:</b> ${item.id} — ${item.name}: $ ${item.price.toFixed(2)}</span>
            <div class="item-actions">
                <button class="btn-edit" onclick='prepareEdit(${JSON.stringify(item)})'>Edit</button>
                <button class="btn-delete" onclick="deleteItem(${item.id})">Delete</button>
            </div> 
        `;
        ulResult.appendChild(li);
    } catch (err) {
        console.error('Error searching item:', err);
        ulResult.innerHTML = `<li style="color:red;">Error searching item</li>`;
    }
}

function cancelEdit() {
    form.reset();
    itemIdInput.value = '';
    submitButton.textContent = 'Add Item';
    cancelButton.style.display = 'none';
}

async function deleteItem(id) {
    if (!confirm('Are you sure you want to delete this item?')) return;
    try {
        const res = await fetch(`${apiUrl}/${id}`, { method: 'DELETE' });
        if (!res.ok) {
            const msg = await res.text();
            alert('Error deleting item: ' + msg);
            return;
        }
        loadItems();
    } catch (err) {
        console.error('Error deleting item:', err);
        alert('Error deleting item (network): ' + err.message);
    }
}

loadItems();
