// Edit mode JavaScript for dorcs documentation site

(function() {
  const basePath = document.body.dataset.basePath || '';
  let isAuthenticated = false;
  let currentFilePath = null;
  let originalContent = '';

  // =====================
  // Authentication
  // =====================

  async function checkAuth() {
    try {
      const response = await fetch(basePath + '/api/edit/check-auth', {
        credentials: 'include'
      });
      const data = await response.json();
      isAuthenticated = data.authenticated;
      updateAuthUI();
      return isAuthenticated;
    } catch (err) {
      console.error('Auth check failed:', err);
      return false;
    }
  }

  async function login(username, password) {
    try {
      const response = await fetch(basePath + '/api/edit/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username, password })
      });
      const data = await response.json();
      if (data.success) {
        isAuthenticated = true;
        updateAuthUI();
        closeLoginModal();
        return true;
      } else {
        showError(data.message || 'Login failed');
        return false;
      }
    } catch (err) {
      showError('Login failed: ' + err.message);
      return false;
    }
  }

  async function logout() {
    try {
      await fetch(basePath + '/api/edit/logout', {
        method: 'POST',
        credentials: 'include'
      });
      isAuthenticated = false;
      updateAuthUI();
      closeEditMode();
    } catch (err) {
      console.error('Logout failed:', err);
    }
  }

  // =====================
  // Login Modal
  // =====================

  let loginModal = null;

  function createLoginModal() {
    if (loginModal) return loginModal;

    const overlay = document.createElement('div');
    overlay.className = 'edit-overlay';
    overlay.id = 'login-overlay';

    const modal = document.createElement('div');
    modal.className = 'edit-modal';
    modal.innerHTML = `
      <div class="edit-modal-header">
        <h2>Login to Edit Mode</h2>
        <button class="edit-modal-close" onclick="this.closest('.edit-overlay').remove()">×</button>
      </div>
      <div class="edit-modal-body">
        <form id="login-form">
          <div class="edit-form-group">
            <label for="login-username">Username</label>
            <input type="text" id="login-username" required autocomplete="username">
          </div>
          <div class="edit-form-group">
            <label for="login-password">Password</label>
            <input type="password" id="login-password" required autocomplete="current-password">
          </div>
          <div id="login-error" class="edit-error" style="display: none;"></div>
          <button type="submit" class="edit-button edit-button-primary">Login</button>
        </form>
      </div>
    `;

    overlay.appendChild(modal);
    document.body.appendChild(overlay);

    const form = modal.querySelector('#login-form');
    form.addEventListener('submit', async (e) => {
      e.preventDefault();
      const username = document.getElementById('login-username').value;
      const password = document.getElementById('login-password').value;
      await login(username, password);
    });

    overlay.addEventListener('click', (e) => {
      if (e.target === overlay) {
        overlay.remove();
      }
    });

    loginModal = overlay;
    return overlay;
  }

  function openLoginModal() {
    const modal = createLoginModal();
    modal.style.display = 'flex';
    document.getElementById('login-username').focus();
  }

  function closeLoginModal() {
    if (loginModal) {
      loginModal.remove();
      loginModal = null;
    }
  }

  // =====================
  // Edit Mode UI
  // =====================

  let editModePanel = null;
  let fileTree = null;
  let editor = null;

  function createEditModePanel() {
    if (editModePanel) return editModePanel;

    const panel = document.createElement('div');
    panel.className = 'edit-mode-panel';
    panel.id = 'edit-mode-panel';
    panel.innerHTML = `
      <div class="edit-mode-header">
        <h3>Edit Mode</h3>
        <button class="edit-mode-close" onclick="window.dorcsEditMode.close()">×</button>
      </div>
      <div class="edit-mode-content">
        <div class="edit-mode-sidebar">
          <div class="edit-file-tree" id="edit-file-tree">
            <div class="edit-file-tree-header">
              <button class="edit-button edit-button-small" onclick="window.dorcsEditMode.refreshFileTree()">Refresh</button>
              <button class="edit-button edit-button-small" onclick="window.dorcsEditMode.createFile()">New File</button>
              <button class="edit-button edit-button-small" onclick="window.dorcsEditMode.createFolder()">New Folder</button>
            </div>
            <div class="edit-file-tree-content" id="edit-file-tree-content">Loading...</div>
          </div>
        </div>
        <div class="edit-mode-main">
          <div class="edit-editor-header">
            <span id="edit-current-file">No file selected</span>
            <div class="edit-editor-actions">
              <button class="edit-button edit-button-small" onclick="window.dorcsEditMode.saveFile()" id="edit-save-btn">Save</button>
              <button class="edit-button edit-button-small" onclick="window.dorcsEditMode.deleteFile()" id="edit-delete-btn" style="display: none;">Delete</button>
            </div>
          </div>
          <textarea id="edit-editor" class="edit-editor" placeholder="Select a file to edit..."></textarea>
        </div>
      </div>
    `;

    document.body.appendChild(panel);
    editModePanel = panel;
    fileTree = document.getElementById('edit-file-tree-content');
    editor = document.getElementById('edit-editor');

    // Load file tree
    refreshFileTree();

    return panel;
  }

  function openEditMode() {
    if (!isAuthenticated) {
      openLoginModal();
      return;
    }

    const panel = createEditModePanel();
    panel.classList.add('active');
  }

  function closeEditMode() {
    if (editModePanel) {
      editModePanel.classList.remove('active');
      currentFilePath = null;
      editor.value = '';
      originalContent = '';
    }
  }

  async function refreshFileTree(path = '.') {
    if (!fileTree) return;

    try {
      const response = await fetch(basePath + '/api/edit/list-files?path=' + encodeURIComponent(path), {
        credentials: 'include'
      });
      if (!response.ok) throw new Error('Failed to load files');
      const data = await response.json();

      fileTree.innerHTML = '';
      renderFileTree(data.files, path);
    } catch (err) {
      fileTree.innerHTML = '<div class="edit-error">Failed to load files: ' + err.message + '</div>';
    }
  }

  function renderFileTree(files, currentPath) {
    files.sort((a, b) => {
      if (a.is_dir !== b.is_dir) return a.is_dir ? -1 : 1;
      return a.name.localeCompare(b.name);
    });

    files.forEach(file => {
      const item = document.createElement('div');
      item.className = 'edit-file-item' + (file.is_dir ? ' edit-file-dir' : '');
      item.innerHTML = `
        <span class="edit-file-icon">${file.is_dir ? '📁' : '📄'}</span>
        <span class="edit-file-name">${escapeHtml(file.name)}</span>
      `;

      if (file.is_dir) {
        item.addEventListener('click', () => {
          refreshFileTree(file.path);
        });
      } else {
        item.addEventListener('click', () => {
          loadFile(file.path);
        });
      }

      fileTree.appendChild(item);
    });
  }

  async function loadFile(path) {
    try {
      const response = await fetch(basePath + '/api/edit/read-file?path=' + encodeURIComponent(path), {
        credentials: 'include'
      });
      if (!response.ok) throw new Error('Failed to load file');
      const data = await response.json();

      currentFilePath = data.path;
      originalContent = data.content;
      editor.value = data.content;
      document.getElementById('edit-current-file').textContent = data.path;
      document.getElementById('edit-delete-btn').style.display = 'inline-block';
    } catch (err) {
      showError('Failed to load file: ' + err.message);
    }
  }

  async function saveFile() {
    if (!currentFilePath) {
      showError('No file selected');
      return;
    }

    const content = editor.value;
    try {
      const response = await fetch(basePath + '/api/edit/save-file', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ path: currentFilePath, content })
      });
      if (!response.ok) throw new Error('Failed to save file');
      const data = await response.json();

      if (data.success) {
        originalContent = content;
        showSuccess('File saved successfully');
        // Reload page to show changes
        setTimeout(() => window.location.reload(), 500);
      }
    } catch (err) {
      showError('Failed to save file: ' + err.message);
    }
  }

  async function createFile() {
    const name = prompt('Enter file name (e.g., new-doc.md):');
    if (!name) return;

    const path = name;
    try {
      const response = await fetch(basePath + '/api/edit/create-file', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ path, content: '', is_dir: false })
      });
      if (!response.ok) throw new Error('Failed to create file');
      refreshFileTree();
      loadFile(path);
    } catch (err) {
      showError('Failed to create file: ' + err.message);
    }
  }

  async function createFolder() {
    const name = prompt('Enter folder name:');
    if (!name) return;

    const path = name;
    try {
      const response = await fetch(basePath + '/api/edit/create-file', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ path, is_dir: true })
      });
      if (!response.ok) throw new Error('Failed to create folder');
      refreshFileTree();
    } catch (err) {
      showError('Failed to create folder: ' + err.message);
    }
  }

  async function deleteFile() {
    if (!currentFilePath) {
      showError('No file selected');
      return;
    }

    if (!confirm('Are you sure you want to delete ' + currentFilePath + '?')) {
      return;
    }

    try {
      const response = await fetch(basePath + '/api/edit/delete-file', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ path: currentFilePath })
      });
      if (!response.ok) throw new Error('Failed to delete file');
      const data = await response.json();

      if (data.success) {
        currentFilePath = null;
        editor.value = '';
        originalContent = '';
        document.getElementById('edit-current-file').textContent = 'No file selected';
        document.getElementById('edit-delete-btn').style.display = 'none';
        refreshFileTree();
        showSuccess('File deleted successfully');
      }
    } catch (err) {
      showError('Failed to delete file: ' + err.message);
    }
  }

  // =====================
  // UI Updates
  // =====================

  function updateAuthUI() {
    // Login button is in footer - show when NOT authenticated
    const loginBtn = document.getElementById('edit-login-btn');
    
    // Edit and Logout buttons are in header - show when authenticated
    const editBtn = document.getElementById('edit-edit-btn');
    const logoutBtn = document.getElementById('edit-logout-btn');

    if (isAuthenticated) {
      // Hide login button in footer
      if (loginBtn) loginBtn.style.display = 'none';
      // Show edit and logout buttons in header
      if (editBtn) editBtn.style.display = 'inline-block';
      if (logoutBtn) logoutBtn.style.display = 'inline-block';
    } else {
      // Show login button in footer
      if (loginBtn) loginBtn.style.display = 'inline-block';
      // Hide edit and logout buttons in header
      if (editBtn) editBtn.style.display = 'none';
      if (logoutBtn) logoutBtn.style.display = 'none';
    }
  }

  function showError(message) {
    const errorDiv = document.getElementById('login-error') || document.createElement('div');
    errorDiv.className = 'edit-error';
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    if (!errorDiv.parentElement) {
      document.body.appendChild(errorDiv);
      setTimeout(() => errorDiv.remove(), 5000);
    }
  }

  function showSuccess(message) {
    const successDiv = document.createElement('div');
    successDiv.className = 'edit-success';
    successDiv.textContent = message;
    document.body.appendChild(successDiv);
    setTimeout(() => successDiv.remove(), 3000);
  }

  function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
  }

  // =====================
  // Expose API
  // =====================

  window.dorcsEditMode = {
    checkAuth,
    login: openLoginModal,
    logout,
    open: openEditMode,
    close: closeEditMode,
    refreshFileTree,
    createFile,
    createFolder,
    saveFile,
    deleteFile,
    loadFile
  };

  // Check auth on load
  checkAuth();
})();

