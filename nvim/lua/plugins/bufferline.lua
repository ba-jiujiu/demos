return {
    "https://github.com/akinsho/bufferline.nvim",
    dependencies = {
        "nvim-tree/nvim-web-devicons"
    },
    opts = {},
    keys = {
        { "<leader>bh", ":BufferLineCyclePrev<CR>", silent = true, noremap = true },
        { "<leader>bl", ":BufferLineCycleNext<CR>", silent = true, noremap = true },
        { "<leader>bp", ":BufferLinePick<CR>", silent = true, noremap = true },
        { "<leader>bc", ":BufferLinePickClose<CR>", silent = true, noremap = true },
        { "<leader>bo", ":BufferLineCloseOther<CR>", silent = true, noremap = true },
        { "<leader>bd", ":bdelete<CR>", silent = true, noremap = true },
    },
    lazy = false
}
