-- Keymaps are automatically loaded on the VeryLazy event
-- Default keymaps that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/keymaps.lua
-- Add any additional keymaps here
local map = vim.keymap.set
-- local util = require("lazy.util") -- 我们不再需要 util.float_cmd

-- 为 go-zero 创建一个专门的前缀 <leader>gz
-- <leader> 是空格键

--【已修改】格式化 API 文件：使用后台静默执行和通知
map("n", "<leader>gzf", function()
  local file = vim.fn.expand("%")
  local cmd = { "goctl", "api", "format", "--dir", file }
  local output = vim.fn.system(cmd)

  if vim.v.shell_error == 0 then
    vim.notify("API formatted successfully: " .. file, vim.log.levels.INFO)
  else
    vim.notify("API format failed:\n" .. output, vim.log.levels.ERROR)
  end
end, { desc = "[G]o-[Z]ero: [F]ormat current API file" })

--【已修改】从 API 生成 Go 代码：同样使用后台静默执行和通知 (推荐方案)
map("n", "<leader>gzg", function()
  local file = vim.fn.expand("%")
  -- 注意：这里的 -dir . 表示在当前 Neovim 的工作目录下生成代码
  -- 请确保您的工作目录是项目根目录（可以使用 :pwd 查看）
  local cmd = { "goctl", "api", "go", "-api", file, "-dir", "." }
  local output = vim.fn.system(cmd)

  if vim.v.shell_error == 0 then
    vim.notify("Go code generated successfully from: " .. file, vim.log.levels.INFO)
  else
    vim.notify("Go code generation failed:\n" .. output, vim.log.levels.ERROR)
  end
end, { desc = "[G]o-[Z]ero: [G]enerate Go code from API" })
