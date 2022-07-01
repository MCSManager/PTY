const {spawn} = require("child_process");
const p = spawn(
    "main.exe",
    [
        "-dir",
        ".",
        "-cmd",
        'cmd.exe',
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
    let text = v.toString()
    if (text.toString().includes("exit0")) {
        return p.stdin.write(JSON.stringify({
            type: 3,
            data: ""
        })+"\n");
    }
    if (text.toString().includes("resize")) {
        const arr = text.split(" ").slice(1)
        console.log("RESIZE WIN:",`${arr[0]} ${arr[1]}`)
        return p.stdin.write(JSON.stringify({
            type: 2,
            data: `${arr[0]},${arr[1]}`
        })+"\n");
    }
     text = JSON.stringify({
        type: 1,
        data: v.toString()
    })
    console.log("[DEBUG] Node >>>>> Go", text)
    p.stdin.write(text + "\n");
});
