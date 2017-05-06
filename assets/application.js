document.addEventListener('DOMContentLoaded', function() {
    /*
    Auto scale textarea
     */
    var minHeight = 100;
    var noteArea = document.getElementById('note-input-area');
    noteArea.oninput = function() {
        noteArea.style.height = null;
        noteArea.style.height = Math.max(noteArea.scrollHeight, minHeight) + 'px';
    };

    /*
    Select text on note link click
     */
    var noteLink = document.getElementById('note-link');
    if (noteLink) noteLink.onclick = function() { selectText(); };
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
    if (copyBtn) copyBtn.onclick = function(e) {
        e.preventDefault();
        selectText();
        document.execCommand('copy');
    };
});
