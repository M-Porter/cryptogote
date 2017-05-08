{{if .noteLink}}
    <div class="note-link-container">
        <div id="note-link" class="note-link">
            <div>{{.noteLink}}</div>
        </div>
        <button class="btn copy-btn" name="button" type="submit" id="copy-btn">Copy Link</button>
    </div>
{{end}}

<form action="/messages" method="POST" accept-charset="UTF-8">
    {{.csrfField}}
    <div class="note-container">
        <textarea
            class="crypto-input"
            id="note-input-area"
            name="Note.Content"
            style="width: 100%;"
            placeholder="Go crazy..."
            autofocus="autofocus"></textarea>
    </div>
    <div class="note-btn-container">
        <button class="btn" type="submit" id="btn">Create</button>
    </div>
</form>
