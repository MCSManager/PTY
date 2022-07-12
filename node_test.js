const { spawn } = require("child_process");

const command = JSON.stringify(["cmd.exe"]);

const p = spawn("./pty.exe", ["-dir", ".", "-cmd", command, "-size", "50,50"], {
    cwd: ".",
    stdio: "pipe",
    windowsHide: true,
});

if (!p.pid) throw new Error("[DEBUG] 启动失败");

p.on("exit", (err) => {
    console.log("[DEBUG] 程序退出：", err);
});

p.stdout.on("data", (v) => {
    process.stdout.write(v);
});

process.stdin.on("data", (v) => {
    let text = v.toString();
    p.stdin.write(text);
});
