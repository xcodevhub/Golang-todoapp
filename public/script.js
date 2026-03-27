const API_BASE = 'http://127.0.0.1:5050/api/v1';

let currentUserId = null;
let currentTaskId = null;
let allUsers = [];
let allTasks = [];

// Settings
const SETTINGS = {
    theme: localStorage.getItem('theme') || 'light',
    compactMode: localStorage.getItem('compactMode') === 'true',
    notificationsEnabled: localStorage.getItem('notificationsEnabled') !== 'false',
    soundsEnabled: localStorage.getItem('soundsEnabled') !== 'false',
    autoSave: localStorage.getItem('autoSave') !== 'false'
};

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    initializeApp();
});

function initializeApp() {
    // Setup theme
    applyTheme();
    setupThemeToggle();
    setupNavigation();
    setupSearch();
    
    // Load initial data
    switchView('dashboard');
    updateCounts();
}

// Theme Management
function setupThemeToggle() {
    const themeSelect = document.getElementById('theme-select-header');
    
    if (themeSelect) {
        themeSelect.value = SETTINGS.theme;
        themeSelect.addEventListener('change', () => {
            setTheme(themeSelect.value);
        });
    }
}

function setTheme(theme) {
    if (theme === 'auto') {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        theme = prefersDark ? 'dark' : 'light';
    }
    
    if (theme === 'dark') {
        document.documentElement.classList.add('dark-mode');
    } else {
        document.documentElement.classList.remove('dark-mode');
    }
    
    SETTINGS.theme = theme;
    localStorage.setItem('theme', theme);
    
    // Update all selectors
    const headerSelect = document.getElementById('theme-select-header');
    if (headerSelect) {
        headerSelect.value = theme === 'auto' ? theme : (document.documentElement.classList.contains('dark-mode') ? 'dark' : 'light');
    }
}

function applyTheme() {
    if (SETTINGS.theme === 'dark' || 
        (SETTINGS.theme === 'auto' && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        document.documentElement.classList.add('dark-mode');
    } else {
        document.documentElement.classList.remove('dark-mode');
    }
    
    // Update header selector
    const headerSelect = document.getElementById('theme-select-header');
    if (headerSelect) {
        headerSelect.value = SETTINGS.theme;
    }
}

function updateThemeSelect() {
    const headerSelect = document.getElementById('theme-select-header');
    if (headerSelect) {
        headerSelect.value = SETTINGS.theme;
    }
}

function changeTheme() {
    const headerSelect = document.getElementById('theme-select-header');
    if (headerSelect) {
        setTheme(headerSelect.value);
    }
}

// Navigation
function setupNavigation() {
    document.querySelectorAll('.nav-item[data-view]').forEach(item => {
        item.addEventListener('click', () => {
            const view = item.getAttribute('data-view');
            switchView(view);
        });
    });
}

function switchView(viewName) {
    // Remove active class
    document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
    document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
    
    // Add active class
    const view = document.getElementById(`${viewName}-view`);
    const navItem = document.querySelector(`[data-view="${viewName}"]`);
    
    if (view) {
        view.classList.add('active');
    }
    if (navItem) {
        navItem.classList.add('active');
    }
    
    // Update breadcrumb
    const breadcrumbMap = {
        'dashboard': 'Панель управління',
        'users': 'Користувачі',
        'tasks': 'Завдання',
        'statistics': 'Статистика',
        'settings': 'Налаштування',
        'help': 'Довідка'
    };
    
    document.getElementById('current-view').textContent = breadcrumbMap[viewName] || viewName;
    
    // Load specific data
    if (viewName === 'dashboard') {
        loadDashboard();
    } else if (viewName === 'users') {
        loadUsers();
    } else if (viewName === 'tasks') {
        loadTasks();
    } else if (viewName === 'statistics') {
        loadStatistics();
    } else if (viewName === 'settings') {
        loadSettings();
    }
}

// Search
function setupSearch() {
    const searchInput = document.getElementById('search-input');
    if (searchInput) {
        searchInput.addEventListener('input', (e) => {
            performGlobalSearch(e.target.value);
        });
    }
}

function performGlobalSearch(query) {
    if (!query) return;
    console.log('Searching for:', query);
    // Implement global search logic
}

// ===== DASHBOARD =====

async function loadDashboard() {
    try {
        // Load stats for dashboard
        await Promise.all([loadUsers(), loadTasks(), loadStatistics()]);
        
        const usersCount = allUsers.length;
        const tasksCount = allTasks.length;
        const inProgressCount = allTasks.filter(t => t.status === 'in_progress').length;
        const completedCount = allTasks.filter(t => t.status === 'completed').length;
        
        document.getElementById('dashboard-users').textContent = usersCount;
        document.getElementById('dashboard-tasks').textContent = tasksCount;
        document.getElementById('dashboard-in-progress').textContent = inProgressCount;
        document.getElementById('dashboard-completed').textContent = completedCount;
        
        // Recent tasks
        const recentTasks = allTasks.slice().reverse().slice(0, 5);
        const recentContainer = document.getElementById('recent-tasks');
        
        if (recentTasks.length === 0) {
            recentContainer.innerHTML = '<p class="text-muted">Немає недавніх завдань</p>';
            return;
        }
        
        recentContainer.innerHTML = recentTasks.map(task => `
            <div class="recent-item" onclick="switchView('tasks')">
                <div class="recent-item-title">${task.title}</div>
                <div class="recent-item-meta">
                    <span class="status-badge ${task.status}">${getStatusText(task.status)}</span>
                </div>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading dashboard:', error);
    }
}

// ===== USERS FUNCTIONALITY =====

function toggleUserForm() {
    const form = document.getElementById('user-form-container');
    const title = document.getElementById('user-form-title');
    
    if (form.classList.contains('hidden')) {
        form.classList.remove('hidden');
        currentUserId = null;
        title.textContent = 'Додати нового користувача';
        document.getElementById('user-form').reset();
    } else {
        form.classList.add('hidden');
    }
}

async function loadUsers() {
    try {
        const response = await fetch(`${API_BASE}/users`);
        if (!response.ok) throw new Error('Failed to load users');
        
        let users = await response.json();
        allUsers = Array.isArray(users) ? users : users.data || users.users || [];
        
        displayUsers();
        updateUserCount();
    } catch (error) {
        showToast(`Помилка завантаження користувачів: ${error.message}`, 'error');
    }
}

function displayUsers() {
    const list = document.getElementById('users-list');
    
    if (allUsers.length === 0) {
        list.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-users"></i>
                <p>Користувачів не знайдено</p>
            </div>
        `;
        return;
    }
    
    // Apply filter
    const filterValue = document.getElementById('users-filter')?.value.toLowerCase() || '';
    const filteredUsers = allUsers.filter(u => 
        u.full_name.toLowerCase().includes(filterValue) || 
        u.email.toLowerCase().includes(filterValue)
    );
    
    list.innerHTML = filteredUsers.map(user => `
        <div class="list-item">
            <div class="list-item-content">
                <div class="list-item-title">${user.full_name}</div>
                <div class="list-item-text"><i class="fas fa-envelope"></i> ${user.email || 'N/A'}</div>
                <div class="list-item-text"><i class="fas fa-id-card"></i> ID: ${user.id}</div>
            </div>
            <div class="list-item-actions">
                <button class="btn btn-secondary btn-sm" onclick="viewUserDetails(${user.id})">
                    <i class="fas fa-eye"></i> Переглянути
                </button>
                <button class="btn btn-warning btn-sm" onclick="editUser(${user.id})">
                    <i class="fas fa-edit"></i> Редагувати
                </button>
                <button class="btn btn-danger btn-sm" onclick="deleteUser(${user.id})">
                    <i class="fas fa-trash"></i> Видалити
                </button>
            </div>
        </div>
    `).join('');
}

async function handleUserSubmit(event) {
    event.preventDefault();
    
    const fullName = document.getElementById('user-fullname').value;
    const email = document.getElementById('user-email').value;
    
    try {
        let response;
        if (currentUserId) {
            response = await fetch(`${API_BASE}/users/${currentUserId}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ full_name: fullName, email: email })
            });
        } else {
            response = await fetch(`${API_BASE}/users`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ full_name: fullName, email: email })
            });
        }
        
        if (!response.ok) throw new Error('Failed to save user');
        
        showToast(currentUserId ? 'Користувач оновлено' : 'Користувача додано', 'success');
        toggleUserForm();
        loadUsers();
    } catch (error) {
        showToast(`Помилка: ${error.message}`, 'error');
    }
}

function editUser(userId) {
    currentUserId = userId;
    const user = allUsers.find(u => u.id === userId);
    if (user) {
        document.getElementById('user-fullname').value = user.full_name;
        document.getElementById('user-email').value = user.email;
        document.getElementById('user-form-title').textContent = 'Редагувати користувача';
        document.getElementById('user-form-container').classList.remove('hidden');
    }
}

async function deleteUser(userId) {
    if (!confirm('Ви впевнені, що хочете видалити цього користувача?')) return;
    
    try {
        const response = await fetch(`${API_BASE}/users/${userId}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Failed to delete user');
        
        showToast('Користувач видалено', 'success');
        loadUsers();
    } catch (error) {
        showToast(`Помилка видалення: ${error.message}`, 'error');
    }
}

async function viewUserDetails(userId) {
    const user = allUsers.find(u => u.id === userId);
    if (!user) return;
    
    const modal = document.getElementById('modal');
    const title = document.getElementById('modal-title');
    const body = document.getElementById('modal-body');
    
    title.textContent = user.full_name;
    body.innerHTML = `
        <div class="modal-details">
            <p><strong>ID:</strong> ${user.id}</p>
            <p><strong>Повне ім'я:</strong> ${user.full_name}</p>
            <p><strong>Email:</strong> ${user.email || 'N/A'}</p>
            <p><strong>Дата створення:</strong> ${new Date(user.created_at).toLocaleString('uk-UA')}</p>
            <p><strong>Останнє оновлення:</strong> ${new Date(user.updated_at).toLocaleString('uk-UA')}</p>
        </div>
    `;
    
    modal.classList.remove('hidden');
}

function updateUserCount() {
    const badge = document.getElementById('users-count');
    if (badge) {
        badge.textContent = allUsers.length;
        badge.style.display = allUsers.length > 0 ? 'block' : 'none';
    }
}

// ===== TASKS FUNCTIONALITY =====

function toggleTaskForm() {
    const form = document.getElementById('task-form-container');
    const title = document.getElementById('task-form-title');
    
    if (form.classList.contains('hidden')) {
        form.classList.remove('hidden');
        currentTaskId = null;
        title.textContent = 'Додати нове завдання';
        document.getElementById('task-form').reset();
    } else {
        form.classList.add('hidden');
    }
}

async function loadTasks() {
    try {
        const response = await fetch(`${API_BASE}/tasks`);
        if (!response.ok) throw new Error('Failed to load tasks');
        
        let tasks = await response.json();
        allTasks = Array.isArray(tasks) ? tasks : tasks.data || tasks.tasks || [];
        
        displayTasks();
        updateTaskCount();
    } catch (error) {
        showToast(`Помилка завантаження завдань: ${error.message}`, 'error');
    }
}

function displayTasks() {
    const list = document.getElementById('tasks-list');
    
    if (allTasks.length === 0) {
        list.innerHTML = `
            <div class="empty-state">
                <i class="fas fa-tasks"></i>
                <p>Завдань не знайдено</p>
            </div>
        `;
        return;
    }
    
    // Apply filters
    const statusFilter = document.getElementById('tasks-status-filter')?.value || '';
    const searchFilter = document.getElementById('tasks-filter')?.value.toLowerCase() || '';
    
    let filtered = allTasks.filter(t => {
        const matchesStatus = !statusFilter || t.status === statusFilter;
        const matchesSearch = t.title.toLowerCase().includes(searchFilter);
        return matchesStatus && matchesSearch;
    });
    
    list.innerHTML = filtered.map(task => `
        <div class="list-item">
            <div class="list-item-content">
                <div class="list-item-title">${task.title}</div>
                <div class="list-item-text">${task.description || 'Без опису'}</div>
                <div class="list-item-text">
                    <span class="status-badge ${task.status}">${getStatusText(task.status)}</span>
                </div>
            </div>
            <div class="list-item-actions">
                <button class="btn btn-secondary btn-sm" onclick="viewTaskDetails(${task.id})">
                    <i class="fas fa-eye"></i> Переглянути
                </button>
                <button class="btn btn-warning btn-sm" onclick="editTask(${task.id})">
                    <i class="fas fa-edit"></i> Редагувати
                </button>
                <button class="btn btn-danger btn-sm" onclick="deleteTask(${task.id})">
                    <i class="fas fa-trash"></i> Видалити
                </button>
            </div>
        </div>
    `).join('');
}

async function handleTaskSubmit(event) {
    event.preventDefault();
    
    const title = document.getElementById('task-title').value;
    const description = document.getElementById('task-description').value;
    const status = document.getElementById('task-status').value;
    
    try {
        let response;
        if (currentTaskId) {
            response = await fetch(`${API_BASE}/tasks/${currentTaskId}`, {
                method: 'PATCH',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ title, description, status })
            });
        } else {
            response = await fetch(`${API_BASE}/tasks`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ title, description, status })
            });
        }
        
        if (!response.ok) throw new Error('Failed to save task');
        
        showToast(currentTaskId ? 'Завдання оновлено' : 'Завдання додано', 'success');
        toggleTaskForm();
        loadTasks();
    } catch (error) {
        showToast(`Помилка: ${error.message}`, 'error');
    }
}

function editTask(taskId) {
    currentTaskId = taskId;
    const task = allTasks.find(t => t.id === taskId);
    if (task) {
        document.getElementById('task-title').value = task.title;
        document.getElementById('task-description').value = task.description;
        document.getElementById('task-status').value = task.status;
        document.getElementById('task-form-title').textContent = 'Редагувати завдання';
        document.getElementById('task-form-container').classList.remove('hidden');
    }
}

async function deleteTask(taskId) {
    if (!confirm('Ви впевнені, що хочете видалити це завдання?')) return;
    
    try {
        const response = await fetch(`${API_BASE}/tasks/${taskId}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Failed to delete task');
        
        showToast('Завдання видалено', 'success');
        loadTasks();
    } catch (error) {
        showToast(`Помилка видалення: ${error.message}`, 'error');
    }
}

async function viewTaskDetails(taskId) {
    const task = allTasks.find(t => t.id === taskId);
    if (!task) return;
    
    const modal = document.getElementById('modal');
    const title = document.getElementById('modal-title');
    const body = document.getElementById('modal-body');
    
    title.textContent = task.title;
    body.innerHTML = `
        <div class="modal-details">
            <p><strong>ID:</strong> ${task.id}</p>
            <p><strong>Назва:</strong> ${task.title}</p>
            <p><strong>Опис:</strong> ${task.description || 'N/A'}</p>
            <p><strong>Статус:</strong> <span class="status-badge ${task.status}">${getStatusText(task.status)}</span></p>
            <p><strong>Дата створення:</strong> ${new Date(task.created_at).toLocaleString('uk-UA')}</p>
            <p><strong>Останнє оновлення:</strong> ${new Date(task.updated_at).toLocaleString('uk-UA')}</p>
        </div>
    `;
    
    modal.classList.remove('hidden');
}

function filterTasks() {
    displayTasks();
}

function sortTasks() {
    const sortBy = document.getElementById('tasks-sort-filter')?.value || 'created';
    
    if (sortBy === 'created') {
        allTasks.reverse();
    } else if (sortBy === 'updated') {
        allTasks.sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at));
    } else if (sortBy === 'title') {
        allTasks.sort((a, b) => a.title.localeCompare(b.title));
    }
    
    displayTasks();
}

function updateTaskCount() {
    const badge = document.getElementById('tasks-count');
    if (badge) {
        badge.textContent = allTasks.length;
        badge.style.display = allTasks.length > 0 ? 'block' : 'none';
    }
}

// ===== STATISTICS =====

async function loadStatistics() {
    try {
        const response = await fetch(`${API_BASE}/statistics`);
        if (!response.ok) throw new Error('Failed to load statistics');
        
        const stats = await response.json();
        
        const container = document.getElementById('statistics');
        
        container.innerHTML = `
            <div class="stat-card primary">
                <div class="stat-icon"><i class="fas fa-users"></i></div>
                <div class="stat-content">
                    <div class="stat-label">Всього користувачів</div>
                    <div class="stat-value">${stats.total_users || 0}</div>
                </div>
            </div>
            
            <div class="stat-card success">
                <div class="stat-icon"><i class="fas fa-list"></i></div>
                <div class="stat-content">
                    <div class="stat-label">Всього завдань</div>
                    <div class="stat-value">${stats.total_tasks || 0}</div>
                </div>
            </div>
            
            <div class="stat-card warning">
                <div class="stat-icon"><i class="fas fa-hourglass-half"></i></div>
                <div class="stat-content">
                    <div class="stat-label">Завдань в процесі</div>
                    <div class="stat-value">${stats.in_progress_tasks || 0}</div>
                </div>
            </div>
            
            <div class="stat-card info">
                <div class="stat-icon"><i class="fas fa-check-circle"></i></div>
                <div class="stat-content">
                    <div class="stat-label">Завершених завдань</div>
                    <div class="stat-value">${stats.completed_tasks || 0}</div>
                </div>
            </div>
            
            <div class="stat-card">
                <div class="stat-icon"><i class="fas fa-percentage"></i></div>
                <div class="stat-content">
                    <div class="stat-label">% Завершено</div>
                    <div class="stat-value">${stats.total_tasks > 0 ? Math.round((stats.completed_tasks / stats.total_tasks) * 100) : 0}%</div>
                </div>
            </div>
        `;
    } catch (error) {
        showToast(`Помилка завантаження статистики: ${error.message}`, 'error');
    }
}

// ===== SETTINGS =====

function loadSettings() {
    // Settings are already loaded
}

function toggleCompactMode() {
    const checkbox = document.getElementById('compact-mode');
    SETTINGS.compactMode = checkbox.checked;
    localStorage.setItem('compactMode', SETTINGS.compactMode);
    // Apply compact mode styles
    document.body.classList.toggle('compact-mode', SETTINGS.compactMode);
}

function toggleNotifications() {
    const checkbox = document.getElementById('notifications-enabled');
    SETTINGS.notificationsEnabled = checkbox.checked;
    localStorage.setItem('notificationsEnabled', SETTINGS.notificationsEnabled);
}

function toggleSounds() {
    const checkbox = document.getElementById('sounds-enabled');
    SETTINGS.soundsEnabled = checkbox.checked;
    localStorage.setItem('soundsEnabled', SETTINGS.soundsEnabled);
}

function clearCache() {
    if (confirm('Ви впевнені, що хочете очистити кеш?')) {
        localStorage.clear();
        showToast('Кеш очищено. Сторінка перезагружається...', 'success');
        setTimeout(() => location.reload(), 1500);
    }
}

// ===== UTILITIES =====

function getStatusText(status) {
    const statusMap = {
        'pending': 'Очікує',
        'in_progress': 'В процесі',
        'completed': 'Завершено'
    };
    return statusMap[status] || status;
}

function closeModal(event) {
    if (event && event.target.id !== 'modal') return;
    document.getElementById('modal').classList.add('hidden');
}

function showToast(message, type = 'info') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast show ${type}`;
    
    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

function updateCounts() {
    updateUserCount();
    updateTaskCount();
}

// Filter handlers
if (document.getElementById('users-filter')) {
    document.getElementById('users-filter').addEventListener('input', () => displayUsers());
}

if (document.getElementById('tasks-filter')) {
    document.getElementById('tasks-filter').addEventListener('input', () => displayTasks());
}

// Close modal on Escape key
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        document.getElementById('modal').classList.add('hidden');
    }
});
