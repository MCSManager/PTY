const { spawn } = require("child_process");
const readline = require("readline");

// process.chdir("../");

// const command = JSON.stringify(['"C:\\Program Files\\Java\\jdk-17.0.2\\bin\\java"', "-jar", "paper-1.18.1-215.jar"]);
// const command = JSON.stringify(["TerrariaServer.exe"]);
const command = JSON.stringify(["bash"]);

const p = spawn(
  "./main",
  ["-dir", ".", "-cmd", command, "-size", "80,80", "-color"],
  {
    cwd: ".",
    stdio: "pipe",
    windowsHide: true,
  }
);

if (!p.pid) throw new Error("[DEBUG] ERR: PID IS NULL");
console.log("Process started!");

p.on("exit", (err) => {
  console.log("[DEBUG] OK:", err);
});

const rl = readline.createInterface({
  input: p.stdout,
  crlfDelay: Infinity,
});

rl.on("line", (line = "") => {
  console.log("FirstLine:", line);
  listen(line);
  rl.removeAllListeners();
});

function listen(line) {
  // const processInfo = JSON.parse(line);
  console.log("PTY SubProcess Info:", line);
  p.stdout.on("data", (v = "") => {
    process.stdout.write(v);
  });

  process.stdin.on("data", (v) => {
    p.stdin.write(v);
  });
}
