const { spawn } = require("child_process");

// const command = JSON.stringify(['"C:\\Program Files\\Java\\jdk-17.0.2\\bin\\java"', "-jar", "paper-1.18.1-215.jar"]);
const command = JSON.stringify(["TerrariaServer.exe"]);

const p = spawn(
    "./pty.exe",
    [
        "-dir",
        "D:\\Workspace\\DemoLocalhost\\GoTTY\\server\\",
        "-cmd",
        command,
        "-size",
        "80,80",
        "-color"
    ],
    {
        cwd: ".",
        stdio: "pipe",
        windowsHide: true,
    }
);

if (!p.pid) throw new Error("[DEBUG] 启动失败");

p.on("exit", (err) => {
    console.log("[DEBUG] 程序退出：", err);
});

p.stdout.on("data", (v) => {
    process.stdout.write(v);
});

process.stdin.on("data", (v) => {
    p.stdin.write(v);
});
