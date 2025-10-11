const apiUrl = '/bills';

const form = document.getElementById('form-item');
const itemIdInput = document.getElementById('item-id');
const embasaInput = document.getElementById('embasa');
const coelbaInput = document.getElementById('coelba');
const createdAtInput = document.getElementById('created-at');
const updatedAtInput = document.getElementById('updated-at');
const submitButton = document.getElementById('btn-submit');
const cancelButton = document.getElementById('btn-cancel');

function nowISO() {
    return new Date().toISOString().slice(0, 16);
}

createdAtInput.value = nowISO();

form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const id = itemIdInput.value;
    const bill = {
        embasa: embasaInput.value,
        coelba: coelbaInput.value,
        created_at: createdAtInput.value || nowISO(),
        updated_at: id ? nowISO() : null  // add updated_at if the person want to edit the item
    };

    if (id) {
        await updateItem(id, bill);
    } else {
        await addItem(bill);
    }
});

// format date: to YYYY-MM-DD or YYYY-MM-DDThh:mm for input fields
function formatDateForInput(dateString, isDateTimeLocal = false) {
    if (!dateString) return '';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return ''; // Invalid date
    
    const pad = num => num.toString().padStart(2, '0');
    const year = date.getFullYear();
    const month = pad(date.getMonth() + 1);
    const day = pad(date.getDate());
    
    if (isDateTimeLocal) {
        const hours = pad(date.getHours());
        const minutes = pad(date.getMinutes());
        return `${year}-${month}-${day}T${hours}:${minutes}`;
    }
    
    return `${year}-${month}-${day}`;
}

// when the page loads, the user sees a default date range from 30 days ago to today
function setDefaultDateRange() {
    const endDate = new Date();
    const startDate = new Date();
    startDate.setDate(endDate.getDate() - 30);
    
    document.getElementById('start-date').value = formatDateForInput(startDate);
    document.getElementById('end-date').value = formatDateForInput(endDate);
}

function clearDateFilter() {
    console.log('Clearing date filters: ...');
    document.getElementById('start-date').value = '';
    document.getElementById('end-date').value = '';
    loadItems();
}

async function applyDateFilter() {
    const startDate = document.getElementById('start-date').value;
    const endDate = document.getElementById('end-date').value;
    
    if (startDate && endDate && new Date(startDate) > new Date(endDate)) {
        alert('The end date must be greater than or equal to the start date.');
        return;
    }
    
    await loadItems();
}

async function loadItems() {
    try {
        const startDate = document.getElementById('start-date').value;
        const endDate = document.getElementById('end-date').value;
        
        console.log('Loading items with filters:', { startDate, endDate });
        
        // build the URL with the date parameters
        const params = new URLSearchParams();
       
        if (startDate) params.append('start', startDate);
        if (endDate) params.append('end', endDate);
        
        const url = `${apiUrl}?${params.toString()}`;
        
        console.log('Fetching from:', url);
        
        const response = await fetch(url);
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Server response:', errorText);
            throw new Error(`Error ${response.status}: ${response.statusText}`);
        }
        
        const items = await response.json();
        const ul = document.getElementById('item-list');
        ul.innerHTML = '';

        if (!items || items.length === 0) {
            ul.innerHTML = '<li>No bills found for the selected date range</li>';
            return;
        }
        
        items.forEach(item => {
            const li = document.createElement('li');
            li.classList.add("bill-card");
            li.innerHTML = `
                <div class="bill-header">
                    <span class="bill-id">${item.id}</span>
                </div>
                <div class="bill-body">
                    <p><b>EMBASA:</b> R$ ${parseFloat(item.embasa).toFixed(2)}</p>
                    <p><b>COELBA:</b> R$ ${parseFloat(item.coelba).toFixed(2)}</p>
                    <p><b>Created:</b> ${formatDate(item.created_at)}</p>
                    <p><b>Updated:</b> ${formatDate(item.updated_at)}</p>
                </div>
                <div class="bill-actions">
                    <button class="btn-edit" onclick='prepareEdit(${JSON.stringify(item)})'>Edit</button>
                    <button class="btn-delete" onclick="deleteItem('${item.id}')">Delete</button>
                </div>
            `;
            ul.appendChild(li);
        });
        
    } catch (err) {
        console.error('Error loading items:', err);
    }
}

// Need convert date string to ISO format for no problems when sending to API
function toISO(dateString) {
    if (!dateString) return new Date().toISOString();
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        console.warn('Invalid date string:', dateString);
        return new Date().toISOString();
    }
    return date.toISOString();
}

// Form submission handler
form.addEventListener('submit', async (event) => {
   
    event.preventDefault();
    
    try {
        const id = itemIdInput.value;
        const isEditing = !!id;
        
        // Validate numeric values
        const embasaValue = parseFloat(embasaInput.value);
        const coelbaValue = parseFloat(coelbaInput.value);
        
        if (isNaN(embasaValue) || isNaN(coelbaValue) || embasaValue < 0 || coelbaValue < 0) {
            throw new Error('Please enter valid positive numbers for Embasa and Coelba!!');
        }

        const now = new Date().toISOString();
        
        const bill = {
            embasa: embasaValue,
            coelba: coelbaValue,
            created_at: isEditing && createdAtInput.value 
                ? new Date(createdAtInput.value).toISOString() 
                : now,
            updated_at: now
        };

        console.log(isEditing ? 'Updating bill:' : 'Creating new bill:', bill);
        
        if (isEditing) {
            await updateItem(id, bill);
        } else {
            await addItem(bill);
        }
        
        form.reset();
        loadItems();
    } catch (error) {
        console.error('Form submission error:', error);
        //alert(`Error: ${error.message || 'Failed to save bill'}`);
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
            console.error('Error adding item:', msg);
            return;
        }
        
        form.reset();
        loadItems();
        return await res.json();
    } catch (err) {
        console.error('Error adding item:', err);
        //alert('Error adding item. Check the console for more details.')
    }
}

function formatDate(dateString) {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString('pt-BR');
}

// format Numbers and Dates for the form
function prepareEdit(item) {
    itemIdInput.value = item.id;
    embasaInput.value = parseFloat(item.embasa).toFixed(2);
    coelbaInput.value = parseFloat(item.coelba).toFixed(2);
    createdAtInput.value = item.created_at ? item.created_at.slice(0, 16) : nowISO();

    const updatedWrapper = document.getElementById('updated-at-wrapper');
    updatedWrapper.style.display = 'block';
    updatedAtInput.value = nowISO();

    submitButton.textContent = 'Update';
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
            return;
        }
        cancelEdit();
        loadItems();
    } catch (err) {
       console.error('Error updating item:', err);
    }
}

/*async function searchItem() {
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
            <span>
                <b>ID:</b> ${item.id} | 
                <b>EMBASA:</b> R$ ${parseFloat(item.embasa).toFixed(2)} | 
                <b>COELBA:</b> R$ ${parseFloat(item.coelba).toFixed(2)} | 
                <b>Created:</b> ${formatDate(item.created_at)} | 
                <b>Updated:</b> ${formatDate(item.updated_at)}
            </span>
            <div class="item-actions">
                <button class="btn-edit" onclick='prepareEdit(${JSON.stringify(item)})'>Edit</button>
                <button class="btn-delete" onclick="deleteItem('${item.id}')">Delete</button>
            </div> 
        `;
        ulResult.appendChild(li);
    } catch (err) {
        console.error('Error searching item:', err);
        ulResult.innerHTML = `<li style="color:red;">Error searching item</li>`;
    }
}*/

function cancelEdit() {
    form.reset();
    itemIdInput.value = '';
    createdAtInput.value = nowISO();

    document.getElementById('updated-at-wrapper').style.display = 'none';
    updatedAtInput.value = '';

    submitButton.textContent = 'Save';
    cancelButton.style.display = 'none';
}

async function deleteItem(id) {
    if (!confirm('Are you sure you want to delete this item???')) return;
    try {
        await fetch(`${apiUrl}/${id}`, { method: 'DELETE' });
        loadItems();
    } catch (err) {

    }
}

loadItems();
