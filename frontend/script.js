const API_URL = 'http://localhost:8080';

// Show/hide sections
function showSection(sectionId) {
    document.querySelectorAll('.section').forEach(section => {
        section.classList.remove('active');
    });
    document.getElementById(sectionId).classList.add('active');
}

// Services
document.getElementById('serviceForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const service = {
        name: document.getElementById('serviceName').value,
        price: parseFloat(document.getElementById('servicePrice').value),
        commission: parseFloat(document.getElementById('serviceCommission').value)
    };

    try {
        const response = await fetch(`${API_URL}/services`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(service)
        });

        if (response.ok) {
            loadServices();
            e.target.reset();
        }
    } catch (error) {
        console.error('Error:', error);
    }
});

async function loadServices() {
    try {
        const response = await fetch(`${API_URL}/services`);
        const services = await response.json();
        
        const servicesList = document.getElementById('servicesList');
        const orderServiceSelect = document.getElementById('orderService');
        
        servicesList.innerHTML = services.map(service => `
            <div class="list-item">
                <div>
                    <strong>${service.name}</strong>
                    <p>Preço: R$ ${service.price.toFixed(2)} | Comissão: ${service.commission}%</p>
                </div>
                <div class="actions">
                    <button class="edit" onclick="editService(${service.id})" title="Editar">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="delete" onclick="deleteService(${service.id})" title="Excluir">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        `).join('');
        
        orderServiceSelect.innerHTML = '<option value="">Selecione o serviço</option>' +
            services.map(service => `
                <option value="${service.id}">${service.name} - R$ ${service.price.toFixed(2)}</option>
            `).join('');
    } catch (error) {
        console.error('Error loading services:', error);
    }
}

// Clients
document.getElementById('clientForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const client = {
        name: document.getElementById('clientName').value,
        balance: parseFloat(document.getElementById('clientBalance').value)
    };

    try {
        const response = await fetch(`${API_URL}/clients`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(client)
        });

        if (response.ok) {
            loadClients();
            e.target.reset();
        }
    } catch (error) {
        console.error('Error:', error);
    }
});

async function loadClients() {
    try {
        const response = await fetch(`${API_URL}/clients`);
        const clients = await response.json();
        
        const clientsList = document.getElementById('clientsList');
        const orderClientSelect = document.getElementById('orderClient');
        
        clientsList.innerHTML = clients.map(client => `
            <div class="list-item">
                <div>
                    <strong>${client.name}</strong>
                    <p>Saldo: R$ ${client.balance.toFixed(2)}</p>
                </div>
                <div class="actions">
                    <button class="edit" onclick="editClient(${client.id})" title="Editar">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="delete" onclick="deleteClient(${client.id})" title="Excluir">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        `).join('');
        
        orderClientSelect.innerHTML = '<option value="">Selecione o cliente</option>' +
            clients.map(client => `
                <option value="${client.id}">${client.name} - Saldo: R$ ${client.balance.toFixed(2)}</option>
            `).join('');
    } catch (error) {
        console.error('Error loading clients:', error);
    }
}

// Orders
document.getElementById('orderForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const order = {
        client_id: parseInt(document.getElementById('orderClient').value),
        service_id: parseInt(document.getElementById('orderService').value)
    };

    try {
        const response = await fetch(`${API_URL}/orders`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(order)
        });

        if (response.ok) {
            loadOrders();
            e.target.reset();
        }
    } catch (error) {
        console.error('Error:', error);
    }
});

async function loadOrders() {
    try {
        const response = await fetch(`${API_URL}/orders`);
        const orders = await response.json();
        
        const ordersList = document.getElementById('ordersList');
        ordersList.innerHTML = orders.map(order => `
            <div class="list-item">
                <div>
                    <strong>${order.client_name}</strong>
                    <p>Serviço: ${order.service_name}</p>
                    <p>Data: ${new Date(order.date).toLocaleString()}</p>
                    <p>Total: R$ ${order.total.toFixed(2)}</p>
                    <span class="status-${order.status.toLowerCase()}">
                        Status: ${order.status}
                    </span>
                </div>
                <div class="actions">
                    ${order.status === 'pending' ? `
                        <button class="edit" onclick="completeOrder(${order.id})" title="Finalizar">
                            <i class="fas fa-check"></i>
                        </button>
                    ` : ''}
                    <button class="delete" onclick="deleteOrder(${order.id})" title="Excluir">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading orders:', error);
    }
}

// Cashflow
document.getElementById('cashflowForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const cashflow = {
        type: document.getElementById('cashflowType').value,
        description: document.getElementById('cashflowDescription').value,
        amount: parseFloat(document.getElementById('cashflowAmount').value)
    };

    try {
        const response = await fetch(`${API_URL}/cashflow`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(cashflow)
        });

        if (response.ok) {
            loadCashFlow();
            e.target.reset();
        }
    } catch (error) {
        console.error('Error:', error);
    }
});

async function loadCashFlow() {
    try {
        const response = await fetch(`${API_URL}/cashflow`);
        const cashflows = await response.json();
        
        const cashflowList = document.getElementById('cashflowList');
        cashflowList.innerHTML = cashflows.map(cf => `
            <div class="list-item">
                <div>
                    <strong>${cf.description}</strong>
                    <p>Data: ${new Date(cf.date).toLocaleString()}</p>
                    <span class="cashflow-${cf.type}">
                        ${cf.type === 'entrada' ? '+' : '-'} R$ ${cf.amount.toFixed(2)}
                    </span>
                </div>
                <div class="actions">
                    <button class="delete" onclick="deleteCashFlow(${cf.id})" title="Excluir">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading cashflow:', error);
    }
}

// Inventory
document.getElementById('inventoryForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const item = {
        item_name: document.getElementById('itemName').value,
        quantity: parseInt(document.getElementById('itemQuantity').value),
        price: parseFloat(document.getElementById('itemPrice').value)
    };

    try {
        const response = await fetch(`${API_URL}/inventory`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(item)
        });

        if (response.ok) {
            loadInventory();
            e.target.reset();
        }
    } catch (error) {
        console.error('Error:', error);
    }
});

async function loadInventory() {
    try {
        const response = await fetch(`${API_URL}/inventory`);
        const items = await response.json();
        
        const inventoryList = document.getElementById('inventoryList');
        inventoryList.innerHTML = items.map(item => `
            <div class="list-item">
                <div>
                    <strong>${item.name}</strong>
                    <p>Quantidade: ${item.quantity}</p>
                    <p>Preço: R$ ${item.price.toFixed(2)}</p>
                </div>
                <div class="actions">
                    <button class="edit" onclick="editInventoryItem(${item.id})" title="Editar">
                        <i class="fas fa-edit"></i>
                    </button>
                    <button class="delete" onclick="deleteInventoryItem(${item.id})" title="Excluir">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading inventory:', error);
    }
}

// Delete handlers
async function deleteService(id) {
    if (confirm('Tem certeza que deseja excluir este serviço?')) {
        try {
            await fetch(`${API_URL}/services/${id}`, { method: 'DELETE' });
            loadServices();
        } catch (error) {
            console.error('Error deleting service:', error);
        }
    }
}

async function deleteClient(id) {
    if (confirm('Tem certeza que deseja excluir este cliente?')) {
        try {
            await fetch(`${API_URL}/clients/${id}`, { method: 'DELETE' });
            loadClients();
        } catch (error) {
            console.error('Error deleting client:', error);
        }
    }
}

async function deleteOrder(id) {
    if (confirm('Tem certeza que deseja excluir esta comanda?')) {
        try {
            await fetch(`${API_URL}/orders/${id}`, { method: 'DELETE' });
            loadOrders();
        } catch (error) {
            console.error('Error deleting order:', error);
        }
    }
}

async function deleteCashFlow(id) {
    if (confirm('Tem certeza que deseja excluir esta movimentação?')) {
        try {
            await fetch(`${API_URL}/cashflow/${id}`, { method: 'DELETE' });
            loadCashFlow();
        } catch (error) {
            console.error('Error deleting cashflow:', error);
        }
    }
}

async function deleteInventoryItem(id) {
    if (confirm('Tem certeza que deseja excluir este item?')) {
        try {
            await fetch(`${API_URL}/inventory/${id}`, { method: 'DELETE' });
            loadInventory();
        } catch (error) {
            console.error('Error deleting inventory item:', error);
        }
    }
}

// Edit handlers
async function editService(id) {
    const name = prompt('Novo nome do serviço:');
    if (name) {
        const price = parseFloat(prompt('Novo preço:'));
        if (!isNaN(price)) {
            const commission = parseFloat(prompt('Nova comissão (%):'));
            if (!isNaN(commission)) {
                try {
                    await fetch(`${API_URL}/services/${id}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ name, price, commission })
                    });
                    loadServices();
                } catch (error) {
                    console.error('Error updating service:', error);
                }
            }
        }
    }
}

async function editClient(id) {
    const name = prompt('Novo nome do cliente:');
    if (name) {
        const balance = parseFloat(prompt('Novo saldo:'));
        if (!isNaN(balance)) {
            try {
                await fetch(`${API_URL}/clients/${id}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name, balance })
                });
                loadClients();
            } catch (error) {
                console.error('Error updating client:', error);
            }
        }
    }
}

async function editInventoryItem(id) {
    const name = prompt('Novo nome do item:');
    if (name) {
        const quantity = parseInt(prompt('Nova quantidade:'));
        if (!isNaN(quantity)) {
            const price = parseFloat(prompt('Novo preço:'));
            if (!isNaN(price)) {
                try {
                    await fetch(`${API_URL}/inventory/${id}`, {
                        method: 'PUT',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ item_name: name, quantity, price })
                    });
                    loadInventory();
                } catch (error) {
                    console.error('Error updating inventory item:', error);
                }
            }
        }
    }
}

// Order status update
async function completeOrder(id) {
    if (confirm('Deseja finalizar esta comanda?')) {
        try {
            await fetch(`${API_URL}/orders/${id}/status`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ status: 'completed' })
            });
            loadOrders();
            loadClients();
            loadCashFlow();
        } catch (error) {
            console.error('Error completing order:', error);
        }
    }
}

// Initial load
document.addEventListener('DOMContentLoaded', () => {
    showSection('services');
    loadServices();
    loadClients();
    loadOrders();
    loadCashFlow();
    loadInventory();
}); 