-- ~/.config/nvim/lua/plugins/rust.lua
return {
  {
    "mrcjkb/rustaceanvim",
    opts = {
      server = {
        settings = {
          ["rust-analyzer"] = {
            -- 在保存时启用 clippy 检查
            checkOnSave = {
              command = "clippy",
            },
          },
        },
      },
    },
  },
}
