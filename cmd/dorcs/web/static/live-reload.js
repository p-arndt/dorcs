// Live reload functionality for dorcs development mode

(function() {
  let eventSource;
  let reconnectAttempts = 0;
  const maxReconnectAttempts = 5;
  let isConnected = false;

  // Get base path from data attribute
  const basePath = document.body.dataset.basePath || '';

  function getReconnectDelay() {
    // Exponential backoff: 1s, 2s, 4s, 8s, 16s
    return Math.min(1000 * Math.pow(2, reconnectAttempts), 16000);
  }

  function connect() {
    try {
      eventSource = new EventSource(basePath + '/__reload');

      eventSource.onopen = function() {
        if (!isConnected) {
          console.log('[dorcs] Live reload enabled');
          isConnected = true;
        }
        reconnectAttempts = 0;
      };

      eventSource.onmessage = function(event) {
        if (event.data === 'reload') {
          console.log('[dorcs] Changes detected, updating content...');
          smartReload();
        } else if (event.data === 'fullreload') {
          console.log('[dorcs] Config changes detected, reloading page...');
          window.location.reload();
        }
      };

      eventSource.onerror = function(e) {
        eventSource.close();

        // Only attempt reconnect if we had a successful connection before
        if (isConnected && reconnectAttempts < maxReconnectAttempts) {
          reconnectAttempts++;
          const delay = getReconnectDelay();
          setTimeout(connect, delay);
        } else if (reconnectAttempts >= maxReconnectAttempts) {
          console.log('[dorcs] Live reload: max reconnection attempts reached');
        }
      };
    } catch (e) {
      console.error('[dorcs] Live reload error:', e);
    }
  }

  // Smart reload: update content without full page refresh
  function smartReload() {
    // Save current scroll position
    const scrollY = window.scrollY;
    const scrollX = window.scrollX;

    // Save sidebar state (open/closed folders)
    const expandedFolders = new Set();
    document.querySelectorAll('.nav-tree .folder:not(.collapsed)').forEach(function(folder) {
      const node = folder.querySelector('.node');
      if (node) {
        const link = node.querySelector('a');
        if (link) {
          expandedFolders.add(link.getAttribute('href'));
        }
      }
    });

    // Fetch the updated page
    fetch(window.location.href, {
      headers: { 'X-Requested-With': 'XMLHttpRequest' }
    })
    .then(function(response) {
      if (!response.ok) {
        throw new Error('Failed to fetch updated content');
      }
      return response.text();
    })
    .then(function(html) {
      const parser = new DOMParser();
      const newDoc = parser.parseFromString(html, 'text/html');

      // Update main content
      const currentMain = document.querySelector('section.main');
      const newMain = newDoc.querySelector('section.main');
      if (currentMain && newMain) {
        currentMain.innerHTML = newMain.innerHTML;
      }

      // Update navigation if it changed
      const currentNav = document.querySelector('.nav-tree');
      const newNav = newDoc.querySelector('.nav-tree');
      if (currentNav && newNav) {
        currentNav.innerHTML = newNav.innerHTML;

        // Restore expanded folders
        expandedFolders.forEach(function(href) {
          const link = document.querySelector('.nav-tree a[href="' + href + '"]');
          if (link) {
            const node = link.closest('.node');
            if (node) {
              const folder = node.parentElement;
              if (folder && folder.classList.contains('folder')) {
                folder.classList.remove('collapsed');
              }
            }
          }
        });

        // Re-highlight current page
        const currentPath = window.location.pathname;
        const activeLink = document.querySelector('.nav-tree a[href="' + currentPath + '"]');
        if (activeLink) {
          document.querySelectorAll('.nav-tree .node').forEach(function(n) {
            n.classList.remove('active');
          });
          const activeNode = activeLink.closest('.node');
          if (activeNode) {
            activeNode.classList.add('active');
          }
        }
      }

      // Update TOC
      const currentToc = document.querySelector('.toc-body');
      const newToc = newDoc.querySelector('.toc-body');
      if (currentToc && newToc) {
        currentToc.innerHTML = newToc.innerHTML;
      }

      // Update mobile TOC
      const currentMobileToc = document.querySelector('.mobile-toc-body');
      const newMobileToc = newDoc.querySelector('.mobile-toc-body');
      if (currentMobileToc && newMobileToc) {
        currentMobileToc.innerHTML = newMobileToc.innerHTML;
      }

      // Update page title
      const newTitle = newDoc.querySelector('title');
      if (newTitle) {
        document.title = newTitle.textContent;
      }

      // Restore scroll position after a short delay to allow content to render
      setTimeout(function() {
        window.scrollTo(scrollX, scrollY);
      }, 10);

      console.log('[dorcs] Content updated');
    })
    .catch(function(error) {
      console.log('[dorcs] Smart reload failed, doing full reload:', error);
      window.location.reload();
    });
  }

  connect();

  // Clean up on page unload
  window.addEventListener('beforeunload', function() {
    if (eventSource) {
      eventSource.close();
    }
  });
})();

