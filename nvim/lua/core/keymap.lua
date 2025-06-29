vim.keymap.set({ "n", "i" }, "<A-z>", "<Cmd>undo<CR>", { 
    silent = true,
    noremap = true })
vim.keymap.set({ "n", "i" }, "<A-x>", "<Cmd>redo<CR>", {
    silent = true,
    noremap = true })

vim.g.mapleader = " "
