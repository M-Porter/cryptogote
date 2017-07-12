document.addEventListener('DOMContentLoaded', function() {
    /*
    Auto scale textarea
     */
    var minHeight = 100;
    var noteArea = document.getElementById('note-input-area');
    if (noteArea) noteArea.addEventListener('input', function() {
        noteArea.style.height = null;
        noteArea.style.height = Math.max(noteArea.scrollHeight, minHeight) + 'px';
    });

    /*
    Select text on note link click
     */
    var noteLink = document.getElementById('note-link');
    if (noteLink) noteLink.addEventListener('click', function () { selectText(); });
    var selectText = function selectText() {
        if (document.selection) {
            var range = document.body.createTextRange();
            range.moveToElementText(noteLink);
            range.select();
        } else if (window.getSelection) {
            var range = document.createRange();
            range.selectNodeContents(noteLink);
            window.getSelection().removeAllRanges();
            window.getSelection().addRange(range);
        }
    };

    /*
    Copy on click
    */
    var copyBtn = document.getElementById('copy-btn');
    if (copyBtn) copyBtn.addEventListener('click', function(e) {
        e.preventDefault();
        selectText();
        document.execCommand('copy');
    });

    /*
    Focus note area on container click
    */
    var noteContainer = document.querySelector('.note-container');
    if (noteContainer) noteContainer.addEventListener('click', function(e) {
        if (noteArea) noteArea.focus();
    });
});
