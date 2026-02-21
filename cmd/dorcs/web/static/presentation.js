(function() {
  'use strict';
  var deck = document.querySelector('.presentation-deck');
  if (!deck) return;
  var slides = deck.querySelectorAll('.presentation-slide');
  var total = slides.length;
  if (total === 0) return;

  function toggleFullscreen() {
    if (!document.fullscreenElement) {
      document.documentElement.requestFullscreen().catch(function() {});
      document.body.classList.add('presentation-fullscreen');
    } else {
      document.exitFullscreen();
      document.body.classList.remove('presentation-fullscreen');
    }
  }

  function goToSlide(index) {
    index = Math.max(0, Math.min(index, total - 1));
    var slide = slides[index];
    if (slide) slide.scrollIntoView({ behavior: 'smooth' });
    window.location.hash = index + 1;
    updateProgress(index);
  }

  function updateProgress(index) {
    var el = document.querySelector('.presentation-progress');
    if (!el) return;
    var slide = slides[index];
    var paginate = slide ? slide.getAttribute('data-paginate') : '';
    var hide = paginate === 'false' || paginate === 'skip';
    el.style.visibility = hide ? 'hidden' : '';
    el.textContent = (index + 1) + ' / ' + total;
  }

  function currentIndex() {
    var y = deck.scrollTop;
    var idx = 0;
    for (var i = 0; i < slides.length; i++) {
      if (slides[i].offsetTop + slides[i].offsetHeight / 2 > y) {
        idx = i;
        break;
      }
      idx = i;
    }
    return idx;
  }

  window.addEventListener('keydown', function(e) {
    if (e.key === 'f' || e.key === 'F') {
      if (!e.ctrlKey && !e.metaKey && !e.altKey) {
        e.preventDefault();
        toggleFullscreen();
      }
      return;
    }
    var idx = currentIndex();
    if (e.key === 'ArrowRight' || e.key === 'ArrowDown' || e.key === ' ') {
      e.preventDefault();
      goToSlide(idx + 1);
    } else if (e.key === 'ArrowLeft' || e.key === 'ArrowUp') {
      e.preventDefault();
      goToSlide(idx - 1);
    } else if (e.key === 'Home') {
      e.preventDefault();
      goToSlide(0);
    } else if (e.key === 'End') {
      e.preventDefault();
      goToSlide(total - 1);
    }
  });

  document.addEventListener('fullscreenchange', function() {
    if (!document.fullscreenElement) {
      document.body.classList.remove('presentation-fullscreen');
    } else {
      document.body.classList.add('presentation-fullscreen');
    }
  });

  deck.addEventListener('scroll', function() {
    updateProgress(currentIndex());
  });

  if (window.location.hash) {
    var n = parseInt(window.location.hash.slice(1), 10);
    if (n >= 1 && n <= total) goToSlide(n - 1);
  } else {
    updateProgress(0);
  }
})();
