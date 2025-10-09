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
    var version = document.getElementById("version-select").value;

    try {
        const res = await fetch(route, {
            method: 'POST',
            headers: {
                'Content-Type': 'text/plain',
                'X-Zig-Version': version
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

async function runCode() {
    execute(runCmd);
}

async function fmtCode() {
    execute(fmtCmd);
}

function toggleTheme() {
    // Determine new value of toggle and set data-theme.
    let newValue = document.documentElement.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', newValue);
    // Update our ace editor with new theme as well.
    let editorTheme = newValue === 'dark' ? 'ace/theme/vibrant_ink' : 'ace/theme/clouds';
    editor.setTheme(editorTheme);
}

/**
 * Our demo code. Setting this up as a variable to keep the HTML clean.
 * In the future, we can add a map of name/code pairs to allow a select
 * list with various code examples. For example this snippet would be "hello world"
 * but we would also allow them to select "zigg zagg" from the examples:
 * 
 * https://ziglang.org/learn/samples/
 * */
const demoCode = `// You can edit this code!
// Click into the editor and start typing.
const std = @import("std");
const builtin = @import("builtin");

pub fn main() void {
    std.debug.print("Hello, {s}! (using Zig version: {f})", .{ "world", builtin.zig_version });
}
`; // Adding trailing space so it matches "format" output.

// On content loaded, set up ace editor and detect dark/light mode and set data-theme appropriately.
document.addEventListener('DOMContentLoaded', function () {
    editor = ace.edit("editor");
    // Set the value of our editor to our demo code.
    // The second param prevents default "select all" behavior.
    editor.setValue(demoCode, -1);
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
        editor.setTheme("ace/theme/vibrant_ink");
        document.documentElement.setAttribute('data-theme', 'dark');
    } else {
        editor.setTheme("ace/theme/clouds");
        document.documentElement.setAttribute('data-theme', 'light');
    }
    // Set editor mode to zig.
    editor.session.setMode("ace/mode/zig");
});
