/**
 *  Global ace editor variable.
 */
var editor;

const runCmd = '/server/run';
const fmtCmd = '/server/fmt';

/**
 * Execute a function.
 * @param {string} route - The function to execute.
 */
async function execute(route) {
    const stdout = document.getElementById('stdout');
    if (!stdout) {
        console.error("Couldn't find element #stdout.");
        return;
    }
    stdout.innerHTML = "Waiting for server...";
    const code = editor.getValue();
    try {
        const res = await fetch(route, {
            method: 'POST',
            headers: {
                'Content-Type': 'text/plain'
            },
            body: code
        });
        const text = await res.text();
        let msg = text;
        if (res.status === 429) {
            msg = 'Too many requests. Please wait a minute and then try again.';
        } else if (!res.ok) {
            msg = 'An error occurred:\n' + text;
        }
        // If run command, we display the resulting text.
        if (route === runCmd) {
            stdout.innerHTML = msg;
        } else if (route === fmtCmd) {
            // For format command, we set the editor text to the output and 
            // let the user know the command is complete.
            if (code != msg) {
                // Preserve selection to try to put the cursor about where
                // it was before format.
                selection = editor.selection.toJSON();
                editor.setValue(msg);
                editor.selection.fromJSON(selection)
            }
            stdout.innerHTML = '';
        }
    } catch (e) {
        stdout.innerHTML = 'Could not connect to server.';
    }
}

/**
 * Executes `runCmd` on the code in the editor.
 */
async function runCode() {
    execute(runCmd);
}

/**
 * Executes `fmtCmd` on the code in the editor then replaces the code in the editor with formatted code.
 */
async function fmtCode() {
    execute(fmtCmd);
}

/**
 * Toggle the dark/light theme.
 */
function toggleTheme() {
    // Determine new value of toggle and set data-theme.
    let newValue = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', newValue);
    // Update our ace editor with new theme as well.
    let editorTheme = newValue === 'dark' ? 'ace/theme/vibrant_ink' : 'ace/theme/clouds';
    editor.setTheme(editorTheme);
}

// On content loaded, set up ace editor and detect dark/light mode and set data-theme appropriately.
document.addEventListener('DOMContentLoaded', function () {
    editor = ace.edit("editor");
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        // dark mode
        editor.setTheme("ace/theme/vibrant_ink");
        document.documentElement.setAttribute('data-theme', 'dark');
    } else {
        // light mode
        editor.setTheme("ace/theme/clouds");
        document.documentElement.setAttribute('data-theme', 'light');
    }
    // Set editor mode to zig.
    editor.session.setMode("ace/mode/zig");
});
