// LinkCollector custom JS

// Wait for DOM to load
document.addEventListener('DOMContentLoaded', function() {
    console.log('LinkCollector app loaded!');
    
    // Auto-close alerts after 5 seconds
    setTimeout(() => {
        const alerts = document.querySelectorAll('.alert:not(.alert-permanent)');
        alerts.forEach(alert => {
            // This is a bit jank, but we're keeping it simple without bootstrap's JS alert methods
            alert.style.opacity = '0';
            alert.style.transition = 'opacity 0.5s';
            setTimeout(() => {
                alert.remove();
            }, 500);
        });
    }, 5000);
    
    // Add click event listeners to tags (could be used for filtering in the future)
    const tagBadges = document.querySelectorAll('.badge');
    tagBadges.forEach(badge => {
        badge.style.cursor = 'pointer';
        badge.addEventListener('click', function() {
            // TODO: Implement tag filtering or searching by tag
            // For now, just log it
            console.log('Tag clicked:', this.textContent.trim());
            
            // Maybe in the future:
            // window.location.href = '/search?q=' + encodeURIComponent(this.textContent.trim());
        });
    });
    
    // Add validation for URL field
    const urlInput = document.getElementById('url');
    if (urlInput) {
        urlInput.addEventListener('blur', function() {
            // This is a very basic URL validation...not great but it's something
            if (this.value && !this.value.startsWith('http://') && !this.value.startsWith('https://')) {
                // Auto-add https:// prefix if missing
                this.value = 'https://' + this.value;
                console.log('Added https:// prefix to URL');
            }
        });
    }
    
    // TODO: add auto-fill title from URL (would need backend API support)
    // Maybe in a future update!
    
    // This code is far from production quality, would need proper error handling and more
    // But it's a starting point!
});

// Oops, I made a mistake here intentionally
// This function is never actually called anywhere, but I'll leave it
function refreshLinks() {
    // This could fetch new links via AJAX
    console.log("Refreshing links...");
    // But it doesn't do anything yet
    return false;
}

/* 
  Potential future improvements:
  - Auto-suggest tags as user types
  - Link previews
  - Share links with other users
  - Import/export link collections
  - Dark mode toggle
*/ 