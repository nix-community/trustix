const ChildProcess = require("child_process")
const Process = require("process")
const Path = require("path")
const FS = require("fs")

const distDir = Path.join(__dirname, "dist")
const srcDir = Path.join(__dirname, "src")

FS.readdir(srcDir, (err, files) => {
  if (err !== null) {
    console.log(err)
    Process.exit(1)
  }

  const filePaths = files.map(f => Path.join(srcDir, f))
  filePaths.forEach(path => {
    console.log(`Building ${path}`)
    const p = ChildProcess.spawnSync("unbundle", [
      "--root", "/js/",
      "--force",
      "--entry", path,
      "--destination", distDir,
    ])
    if (p.status != 0) {
      console.log(`Building ${path} failed with status ${p.status}:`)
      console.log(p.stderr.toString())
      Process.exit(p.status)
    }
  })

})
