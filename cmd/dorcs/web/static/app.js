// Main application JavaScript for dorcs documentation site

// Define switcher functions immediately (before IIFE) so they're available for inline handlers
window.dorcsVersionSwitcher = {
  // Save version preference to localStorage
  saveVersionPreference: function(versionID) {
    try {
      localStorage.setItem('dorcs-version', versionID || '');
    } catch (e) {
      // Ignore localStorage errors (e.g., in private browsing)
    }
  },
  
  // Get saved version preference from localStorage
  getVersionPreference: function() {
    try {
      return localStorage.getItem('dorcs-version') || '';
    } catch (e) {
      return '';
    }
  },
  
  switchVersion: function(versionID, defaultVersion) {
    const basePath = document.body.dataset.basePath || '';
    // Use DocPath from server (already calculated correctly, without version/language prefix)
    let docPath = document.body.dataset.docPath || '/';
    
    // Ensure docPath starts with /
    if (!docPath.startsWith('/')) {
      docPath = '/' + docPath;
    }
    
    // Get current language (if any)
    const currentLang = document.body.dataset.currentLanguage || '';
    const langSelect = document.getElementById('lang-select');
    const defaultLang = langSelect ? langSelect.dataset.defaultLang : '';
    
    // Build new path with the selected version (language-first structure)
    let newPath = '';
    if (versionID === defaultVersion || versionID === '') {
      // Default version - no version prefix
      if (currentLang && currentLang !== defaultLang) {
        // Has language, add language prefix only
        if (docPath === '/') {
          newPath = '/' + currentLang + '/';
        } else {
          newPath = '/' + currentLang + docPath;
        }
      } else {
        // No language or default language
        newPath = docPath;
      }
    } else {
      // Non-default version - add version prefix (after language if present)
      if (currentLang && currentLang !== defaultLang) {
        // Has language: /en/v1/... (language-first)
        if (docPath === '/') {
          newPath = '/' + currentLang + '/' + versionID + '/';
        } else {
          newPath = '/' + currentLang + '/' + versionID + docPath;
        }
      } else {
        // No language or default language: /v1/...
        if (docPath === '/') {
          newPath = '/' + versionID + '/';
        } else {
          newPath = '/' + versionID + docPath;
        }
      }
    }
    
    // Save preference and navigate
    this.saveVersionPreference(versionID);
    window.location.href = basePath + newPath;
  }
};

window.dorcsLangSwitcher = {
  // Save language preference to localStorage
  saveLanguagePreference: function(langCode) {
    try {
      localStorage.setItem('dorcs-language', langCode || '');
    } catch (e) {
      // Ignore localStorage errors (e.g., in private browsing)
    }
  },
  
  // Get saved language preference from localStorage
  getLanguagePreference: function() {
    try {
      return localStorage.getItem('dorcs-language') || '';
    } catch (e) {
      return '';
    }
  },
  
  switchLanguage: function(langCode, defaultLang) {
    const basePath = document.body.dataset.basePath || '';
    // Use DocPath from server (already calculated correctly, without version/language prefix)
    let docPath = document.body.dataset.docPath || '/';
    
    // Ensure docPath starts with /
    if (!docPath.startsWith('/')) {
      docPath = '/' + docPath;
    }
    
    // Get current version (if any)
    const currentVersion = document.body.dataset.currentVersion || '';
    const versionSelect = document.getElementById('version-select');
    const defaultVersion = versionSelect ? versionSelect.dataset.defaultVersion : '';
    
    // Build new path with the selected language (language-first structure)
    let newPath = '';
    if (langCode === defaultLang || langCode === '') {
      // Default language - no language prefix
      if (currentVersion && currentVersion !== defaultVersion) {
        // Has version: /v1/...
        if (docPath === '/') {
          newPath = '/' + currentVersion + '/';
        } else {
          newPath = '/' + currentVersion + docPath;
        }
      } else {
        // No version or default version
        newPath = docPath;
      }
    } else {
      // Other language - add language prefix (before version if present)
      if (currentVersion && currentVersion !== defaultVersion) {
        // Has version: /en/v1/... (language-first)
        if (docPath === '/') {
          newPath = '/' + langCode + '/' + currentVersion + '/';
        } else {
          newPath = '/' + langCode + '/' + currentVersion + docPath;
        }
      } else {
        // No version or default version: /en/...
        if (docPath === '/') {
          newPath = '/' + langCode + '/';
        } else {
          newPath = '/' + langCode + docPath;
        }
      }
    }
    
    // Save preference and navigate
    this.saveLanguagePreference(langCode);
    window.location.href = basePath + newPath;
  }
};

(function() {
  // =====================
  // Mobile Menu Toggle
  // =====================
  const menuToggle = document.getElementById('menu-toggle');
  const sidebar = document.getElementById('sidebar');
  const sidebarOverlay = document.getElementById('sidebar-overlay');

  function closeSidebar() {
    document.body.classList.remove('sidebar-open');
    sidebarOverlay.classList.remove('active');
  }

  function openSidebar() {
    document.body.classList.add('sidebar-open');
    // Small delay for overlay animation
    requestAnimationFrame(() => {
      sidebarOverlay.classList.add('active');
    });
  }

  if (menuToggle) {
    menuToggle.addEventListener('click', () => {
      if (document.body.classList.contains('sidebar-open')) {
        closeSidebar();
      } else {
        openSidebar();
      }
    });
  }

  if (sidebarOverlay) {
    sidebarOverlay.addEventListener('click', closeSidebar);
  }

  // Close sidebar when clicking a link (mobile)
  if (sidebar) {
    sidebar.addEventListener('click', (e) => {
      if (e.target.closest('a[href]') && window.innerWidth <= 768) {
        closeSidebar();
      }
    });
  }

  // Close sidebar on escape key
  document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && document.body.classList.contains('sidebar-open')) {
      closeSidebar();
    }
  });

  // Close sidebar on window resize to desktop
  let resizeTimer;
  window.addEventListener('resize', () => {
    clearTimeout(resizeTimer);
    resizeTimer = setTimeout(() => {
      if (window.innerWidth > 768 && document.body.classList.contains('sidebar-open')) {
        closeSidebar();
      }
    }, 100);
  });

  // =====================
  // Mobile TOC Toggle
  // =====================
  const mobileToc = document.getElementById('mobile-toc');
  const mobileTocToggle = document.getElementById('mobile-toc-toggle');

  if (mobileTocToggle && mobileToc) {
    mobileTocToggle.addEventListener('click', () => {
      mobileToc.classList.toggle('expanded');
    });

    // Close TOC when clicking a link
    mobileToc.addEventListener('click', (e) => {
      if (e.target.closest('a[href^="#"]')) {
        mobileToc.classList.remove('expanded');
      }
    });
  }

  // =====================
  // Collapsible Sidebar
  // =====================
  const navTree = document.getElementById('nav-tree');
  if (navTree) {
    const expandAll = navTree.dataset.expandAll === 'true';
    // Load saved state from localStorage
    const savedState = JSON.parse(localStorage.getItem('nav-collapsed') || '{}');

    // Apply default collapsed state from config.
    if (!expandAll) {
      navTree.querySelectorAll('.folder').forEach(folder => {
        folder.classList.add('collapsed');
      });
    }

    // Apply saved collapsed state
    navTree.querySelectorAll('.folder').forEach(folder => {
      const key = folder.dataset.key;
      if (key && Object.prototype.hasOwnProperty.call(savedState, key)) {
        folder.classList.toggle('collapsed', !!savedState[key]);
      }
    });

    // Auto-expand folders containing the active page
    const activeNode = navTree.querySelector('.node.active');
    if (activeNode) {
      let parent = activeNode.closest('.folder');
      while (parent) {
        parent.classList.remove('collapsed');
        parent = parent.parentElement.closest('.folder');
      }
    }

    // Toggle folders on click
    navTree.addEventListener('click', (e) => {
      const toggle = e.target.closest('.toggle');
      if (toggle) {
        e.preventDefault();
        e.stopPropagation();
        const folder = toggle.closest('.folder');
        if (folder) {
          folder.classList.toggle('collapsed');

          // Save state
          const key = folder.dataset.key;
          if (key) {
            const state = JSON.parse(localStorage.getItem('nav-collapsed') || '{}');
            state[key] = folder.classList.contains('collapsed');
            localStorage.setItem('nav-collapsed', JSON.stringify(state));
          }
        }
      }
    });
  }

  // =====================
  // Search Overlay Modal
  // =====================
  let searchOverlay = null;
  let searchModal = null;
  let searchModalInput = null;
  let searchResultsContainer = null;
  let searchCount = null;
  let searchTimeout = null;
  let currentSearchController = null;
  let selectedIndex = -1;

  // Create search overlay modal
  function createSearchModal() {
    if (searchOverlay) return;

    searchOverlay = document.createElement('div');
    searchOverlay.className = 'search-overlay';
    searchOverlay.id = 'search-overlay';

    searchModal = document.createElement('div');
    searchModal.className = 'search-modal';

    const searchHeader = document.createElement('div');
    searchHeader.className = 'search-header';

    const searchInputWrapper = document.createElement('div');
    searchInputWrapper.className = 'search-input-wrapper';

    const searchIcon = document.createElement('div');
    searchIcon.className = 'search-icon';
    searchIcon.innerHTML = `
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"></circle>
        <path d="m21 21-4.35-4.35"></path>
      </svg>
    `;

    searchModalInput = document.createElement('input');
    searchModalInput.type = 'text';
    searchModalInput.className = 'search-modal-input';
    searchModalInput.placeholder = 'Search documentation...';
    searchModalInput.autocomplete = 'off';

    const closeButton = document.createElement('button');
    closeButton.className = 'search-close';
    closeButton.innerHTML = `
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    `;
    closeButton.setAttribute('aria-label', 'Close search');

    searchInputWrapper.appendChild(searchIcon);
    searchInputWrapper.appendChild(searchModalInput);
    searchInputWrapper.appendChild(closeButton);
    searchHeader.appendChild(searchInputWrapper);

    searchCount = document.createElement('div');
    searchCount.className = 'search-count';

    searchResultsContainer = document.createElement('div');
    searchResultsContainer.className = 'search-results-container';

    searchModal.appendChild(searchHeader);
    searchModal.appendChild(searchCount);
    searchModal.appendChild(searchResultsContainer);
    searchOverlay.appendChild(searchModal);
    document.body.appendChild(searchOverlay);

    // Event listeners
    closeButton.addEventListener('click', closeSearchModal);
    searchOverlay.addEventListener('click', (e) => {
      if (e.target === searchOverlay) {
        closeSearchModal();
      }
    });

    searchModalInput.addEventListener('input', (e) => {
      const query = e.target.value.trim();
      selectedIndex = -1;

      if (searchTimeout) {
        clearTimeout(searchTimeout);
      }

      if (!query) {
        clearSearchResults();
        return;
      }

      searchTimeout = setTimeout(() => {
        performSearch(query);
      }, 300);
    });

    // Keyboard navigation
    searchModalInput.addEventListener('keydown', (e) => {
      const items = searchResultsContainer.querySelectorAll('.search-result-item');
      if (items.length === 0) return;

      if (e.key === 'ArrowDown') {
        e.preventDefault();
        selectedIndex = Math.min(selectedIndex + 1, items.length - 1);
        updateSelection(items);
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        selectedIndex = Math.max(selectedIndex - 1, -1);
        updateSelection(items);
      } else if (e.key === 'Enter' && selectedIndex >= 0) {
        e.preventDefault();
        const selectedItem = items[selectedIndex];
        if (selectedItem) {
          const href = selectedItem.dataset.href;
          if (href) {
            closeSearchModal();
            window.location.href = href;
          }
        }
      } else if (e.key === 'Escape') {
        e.preventDefault();
        closeSearchModal();
      }
    });
  }

  function updateSelection(items) {
    items.forEach((item, i) => {
      item.classList.toggle('selected', i === selectedIndex);
    });
    if (selectedIndex >= 0 && selectedIndex < items.length) {
      items[selectedIndex].scrollIntoView({ block: 'nearest', behavior: 'smooth' });
    }
  }

  function openSearchModal() {
    createSearchModal();
    document.body.classList.add('search-open');
    searchOverlay.classList.add('active');
    setTimeout(() => {
      searchModalInput.focus();
    }, 100);
  }

  function closeSearchModal() {
    if (searchOverlay) {
      document.body.classList.remove('search-open');
      searchOverlay.classList.remove('active');
      searchModalInput.value = '';
      clearSearchResults();
      selectedIndex = -1;
    }
  }

  function clearSearchResults() {
    if (searchResultsContainer) {
      searchResultsContainer.innerHTML = '';
      searchCount.textContent = '';
    }
  }

  function performSearch(query) {
    // Cancel previous request if any
    if (currentSearchController) {
      currentSearchController.abort();
    }

    if (!query || query.trim().length < 1) {
      clearSearchResults();
      return;
    }

    // Get base path from body data attribute
    const basePath = document.body.dataset.basePath || '';
    const searchUrl = (basePath ? basePath : '') + '/api/search?q=' + encodeURIComponent(query.trim());

    // Create new AbortController for this request
    currentSearchController = new AbortController();

    fetch(searchUrl, {
      signal: currentSearchController.signal,
      headers: {
        'Accept': 'application/json'
      }
    })
      .then(response => {
        if (!response.ok) {
          throw new Error('Search failed');
        }
        return response.json();
      })
      .then(data => {
        displaySearchResults(data.results, query);
      })
      .catch(err => {
        if (err.name !== 'AbortError') {
          console.error('Search error:', err);
          displaySearchResults([], query);
        }
      });
  }

  function displaySearchResults(results, query) {
    if (!searchResultsContainer || !searchCount) return;

    const queryLower = query.toLowerCase();
    
    // Escape HTML to prevent XSS and broken rendering
    const escapeHtml = (text) => {
      if (!text) return '';
      const div = document.createElement('div');
      div.textContent = text;
      return div.innerHTML;
    };
    
    const highlightText = (text) => {
      if (!text) return '';
      // First escape HTML
      const escaped = escapeHtml(text);
      // Then apply highlighting
      const regex = new RegExp(`(${queryLower.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi');
      return escaped.replace(regex, '<mark>$1</mark>');
    };

    // Update count
    if (results.length === 0) {
      searchCount.textContent = 'No results found';
      searchResultsContainer.innerHTML = '';
      return;
    }

    const resultText = results.length === 100 ? '100+ matching documents' : `${results.length} matching document${results.length !== 1 ? 's' : ''}`;
    searchCount.textContent = resultText;

    // Build results HTML
    let html = '';
    results.forEach((result, index) => {
      // Escape path to prevent XSS
      const safePath = escapeHtml(result.path || '');
      const safeTitle = highlightText(result.title || '');
      const safeHeading = result.heading_text ? highlightText(result.heading_text) : '';
      const safeSnippet = result.snippet ? highlightText(result.snippet) : '';
      
      const headingInfo = safeHeading ? `<div class="search-result-heading">${safeHeading}</div>` : '';
      html += `
        <div class="search-result-item" data-href="${safePath}" data-index="${index}">
          <div class="search-result-icon">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"></path>
              <polyline points="14 2 14 8 20 8"></polyline>
              <line x1="16" y1="13" x2="8" y2="13"></line>
              <line x1="16" y1="17" x2="8" y2="17"></line>
              <polyline points="10 9 9 9 8 9"></polyline>
            </svg>
          </div>
          <div class="search-result-content">
            <div class="search-result-title">${safeTitle}</div>
            ${headingInfo}
            ${safeSnippet ? `<div class="search-result-snippet">${safeSnippet}</div>` : ''}
            <div class="search-result-path">${safePath}</div>
          </div>
        </div>
      `;
    });
    searchResultsContainer.innerHTML = html;

    // Add click handlers - use mousedown to ensure it fires before blur
    searchResultsContainer.querySelectorAll('.search-result-item').forEach(item => {
      item.addEventListener('mousedown', (e) => {
        e.preventDefault(); // Prevent input blur
        const href = item.dataset.href;
        if (href) {
          closeSearchModal();
          window.location.href = href;
        }
      });
      // Also handle click as fallback
      item.addEventListener('click', (e) => {
        e.preventDefault();
        const href = item.dataset.href;
        if (href) {
          closeSearchModal();
          window.location.href = href;
        }
      });
    });
  }

  // Keyboard shortcut: Ctrl+K or Cmd+K to open search
  document.addEventListener('keydown', (e) => {
    // Don't trigger if user is typing in an input/textarea (unless it's the search modal input)
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA') {
      // Allow Ctrl+K even in inputs to open search
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault();
        openSearchModal();
      }
      return;
    }

    if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
      e.preventDefault();
      openSearchModal();
    }
  });

  // Search toggle button in header
  const searchToggle = document.getElementById('search-toggle');
  if (searchToggle) {
    searchToggle.addEventListener('click', (e) => {
      e.preventDefault();
      openSearchModal();
    });
  }

  // =====================
  // TOC Scrollspy
  // =====================
const tocRoot = document.getElementById('toc');
if (!tocRoot) return;

const HEADER_OFFSET = 80;
const BOTTOM_THRESHOLD = 10;
const SCROLL_EPSILON = 100;

// --- Helpers ---------------------------------------------------------------

function getIdFromLink(link) {
  const href = link.getAttribute('href');
  if (!href || !href.startsWith('#')) return null;
  return decodeURIComponent(href.slice(1));
}

function isScrollable() {
  const { scrollHeight } = document.documentElement;
  return scrollHeight > window.innerHeight + SCROLL_EPSILON;
}

function isAtBottom() {
  const { scrollHeight } = document.documentElement;
  return window.scrollY + window.innerHeight >= scrollHeight - BOTTOM_THRESHOLD;
}

// --- Collect TOC links and targets ----------------------------------------

const links = Array.from(tocRoot.querySelectorAll('a[href^="#"]'));
if (!links.length) return;

const targets = links
  .map(link => {
    const id = getIdFromLink(link);
    if (!id) return null;

    const el = document.getElementById(id);
    return el ? { id, el } : null;
  })
  .filter(Boolean);

if (!targets.length) return;

// Map id → link for fast lookup
const linkById = new Map(
  links
    .map(link => [getIdFromLink(link), link])
    .filter(([id]) => id)
);

// --- Active state management ----------------------------------------------

let currentActiveId = null;

function setActive(id) {
  if (id === currentActiveId) return;

  currentActiveId = id;
  links.forEach(link => link.classList.remove('toc-active'));

  const activeLink = linkById.get(id);
  if (activeLink) activeLink.classList.add('toc-active');
}

// --- Scroll logic ----------------------------------------------------------

function updateActiveSection() {
  const activationLine =
    window.scrollY +
    HEADER_OFFSET +
    window.innerHeight * 0.3; // ← key change

  let activeId = targets[0].id;

  for (const { id, el } of targets) {
    if (el.offsetTop <= activationLine) {
      activeId = id;
    } else {
      break;
    }
  }

  if (isScrollable() && isAtBottom()) {
    setActive(targets[targets.length - 1].id);
    return;
  }

  setActive(activeId);
}


// --- Scroll listener (rAF throttled) ---------------------------------------

let ticking = false;

window.addEventListener(
  'scroll',
  () => {
    if (ticking) return;

    ticking = true;
    requestAnimationFrame(() => {
      updateActiveSection();
      ticking = false;
    });
  },
  { passive: true }
);

// Initial state
updateActiveSection();

  // =====================
  // Code Block Copy Buttons
  // =====================
  function initCopyButtons() {
    document.querySelectorAll('.markdown pre').forEach((pre) => {
      if (pre.querySelector('.copy-button')) return;

      const button = document.createElement('button');
      button.className = 'copy-button';
      button.setAttribute('aria-label', 'Copy code');
      button.innerHTML = `
        <svg class="copy-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
        </svg>
        <svg class="check-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
      `;

      button.addEventListener('click', async () => {
        const code = pre.querySelector('code') || pre;
        const text = code.textContent || code.innerText;
        
        try {
          await navigator.clipboard.writeText(text);
        } catch {
          // Fallback for older browsers
          const textArea = document.createElement('textarea');
          textArea.value = text;
          textArea.style.position = 'fixed';
          textArea.style.opacity = '0';
          document.body.appendChild(textArea);
          textArea.select();
          document.execCommand('copy');
          document.body.removeChild(textArea);
        }
        
        button.classList.add('copied');
        setTimeout(() => button.classList.remove('copied'), 2000);
      });

      pre.style.position = 'relative';
      pre.appendChild(button);
    });
  }

  // Initialize copy buttons on page load
  initCopyButtons();

  // Expose function globally for live-reload to use
  window.dorcsInitCopyButtons = initCopyButtons;

  // =====================
  // Dark Mode Toggle
  // =====================
  (function() {
    const darkModeToggle = document.getElementById('dark-mode-toggle');
    if (!darkModeToggle) return;

    // Apply theme by adding/removing classes
    function applyTheme(theme) {
      const html = document.documentElement;
      html.classList.remove('dark-mode', 'light-mode');
      
      if (theme === 'dark') {
        html.classList.add('dark-mode');
      } else if (theme === 'light') {
        html.classList.add('light-mode');
      }
      // If theme is null/undefined, remove both classes to use server default
      
      if (theme) {
        localStorage.setItem('dorcs-theme', theme);
      } else {
        localStorage.removeItem('dorcs-theme');
      }
    }

    // Get current effective theme (checking classes and system preference)
    function getEffectiveTheme() {
      const html = document.documentElement;
      if (html.classList.contains('dark-mode')) {
        return 'dark';
      }
      if (html.classList.contains('light-mode')) {
        return 'light';
      }
      // No manual override - check system preference
      return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    }

    // Initialize theme on page load (theme may already be applied by inline script in head)
    // This ensures the toggle button state is correct
    const saved = localStorage.getItem('dorcs-theme');
    if (saved === 'dark' || saved === 'light') {
      // Ensure theme is applied (may already be done by inline script, but ensure consistency)
      const html = document.documentElement;
      if (!html.classList.contains(saved + '-mode')) {
        applyTheme(saved);
      }
    }
    // Otherwise, let the server-generated CSS handle it (respects theme.mode config)

    // Toggle theme on button click
    darkModeToggle.addEventListener('click', () => {
      const current = getEffectiveTheme();
      const newTheme = current === 'dark' ? 'light' : 'dark';
      applyTheme(newTheme);
    });
  })();

  // Version and Language Switchers are now defined before the IIFE
  
  // On page load, check if we should redirect to saved language preference
  (function() {
    const savedLang = window.dorcsLangSwitcher.getLanguagePreference();
    const currentLang = document.body.dataset.currentLanguage || '';
    const langSelect = document.getElementById('lang-select');
    const defaultLang = langSelect ? langSelect.dataset.defaultLang : '';
    
    // Only redirect if:
    // 1. We have a saved language preference
    // 2. It's different from current language
    // 3. We're on the default language (root path)
    // 4. The saved language is not the default
    if (savedLang && savedLang !== '' && savedLang !== defaultLang && 
        (!currentLang || currentLang === '') && 
        (window.location.pathname === '/' || window.location.pathname.endsWith('/'))) {
      const basePath = document.body.dataset.basePath || '';
      let pathWithoutBase = window.location.pathname;
      if (basePath && pathWithoutBase.startsWith(basePath)) {
        pathWithoutBase = pathWithoutBase.slice(basePath.length);
      }
      if (pathWithoutBase === '/' || pathWithoutBase === '') {
        // Redirect to saved language
        const newPath = basePath + '/' + savedLang + '/';
        window.location.href = newPath;
      }
    }
  })();

  // === Announcement banner dismissal ===
  (function() {
    var banner = document.getElementById('announcement-banner');
    if (banner) {
      try {
        var dismissed = localStorage.getItem('dorcs-announcement-dismissed');
        if (dismissed && dismissed === banner.querySelector('span').textContent.trim()) {
          banner.remove();
          document.body.classList.remove('has-announcement');
        }
      } catch (e) {}
    }
  })();

  // === Code block copy button ===
  (function() {
    document.querySelectorAll('.markdown pre').forEach(function(pre) {
      // Skip mermaid blocks
      if (pre.classList.contains('mermaid')) return;

      var wrapper = document.createElement('div');
      wrapper.className = 'code-block-wrapper';
      pre.parentNode.insertBefore(wrapper, pre);
      wrapper.appendChild(pre);

      var btn = document.createElement('button');
      btn.className = 'code-copy-btn';
      btn.setAttribute('aria-label', 'Copy code');
      btn.innerHTML = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>';
      btn.addEventListener('click', function() {
        var code = pre.querySelector('code');
        var text = (code || pre).textContent;
        navigator.clipboard.writeText(text).then(function() {
          btn.classList.add('copied');
          btn.innerHTML = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg>';
          setTimeout(function() {
            btn.classList.remove('copied');
            btn.innerHTML = '<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2"/><path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/></svg>';
          }, 2000);
        });
      });
      wrapper.appendChild(btn);
    });
  })();

  // === Heading anchor links ===
  (function() {
    document.querySelectorAll('.markdown h1[id], .markdown h2[id], .markdown h3[id], .markdown h4[id], .markdown h5[id], .markdown h6[id]').forEach(function(heading) {
      var anchor = document.createElement('a');
      anchor.className = 'heading-anchor';
      anchor.href = '#' + heading.id;
      anchor.setAttribute('aria-label', 'Link to this section');
      anchor.textContent = '#';
      anchor.addEventListener('click', function(e) {
        e.preventDefault();
        navigator.clipboard.writeText(window.location.origin + window.location.pathname + '#' + heading.id);
        history.pushState(null, '', '#' + heading.id);
      });
      heading.appendChild(anchor);
    });
  })();

  // === Tabs component ===
  (function() {
    document.querySelectorAll('.dorcs-tabs').forEach(function(tabGroup) {
      var buttons = tabGroup.querySelectorAll('.dorcs-tab-btn');
      var panels = tabGroup.querySelectorAll('.dorcs-tab-panel');
      buttons.forEach(function(btn) {
        btn.addEventListener('click', function() {
          var tabIdx = btn.getAttribute('data-tab');
          buttons.forEach(function(b) { b.classList.remove('active'); });
          panels.forEach(function(p) { p.classList.remove('active'); });
          btn.classList.add('active');
          tabGroup.querySelector('.dorcs-tab-panel[data-tab="' + tabIdx + '"]').classList.add('active');
        });
      });
    });
  })();

  // === Back to top button ===
  (function() {
    var btn = document.createElement('button');
    btn.className = 'back-to-top';
    btn.setAttribute('aria-label', 'Back to top');
    btn.innerHTML = '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="18 15 12 9 6 15"/></svg>';
    btn.addEventListener('click', function() {
      window.scrollTo({ top: 0, behavior: 'smooth' });
    });
    document.body.appendChild(btn);

    var ticking = false;
    window.addEventListener('scroll', function() {
      if (!ticking) {
        requestAnimationFrame(function() {
          if (window.scrollY > 400) {
            btn.classList.add('visible');
          } else {
            btn.classList.remove('visible');
          }
          ticking = false;
        });
        ticking = true;
      }
    });
  })();
})();
