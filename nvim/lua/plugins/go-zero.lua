-- 在 lua/plugins/go-zero.lua 文件中
return {
  {
    "nvim-treesitter/nvim-treesitter",
    opts = function(_, opts)
      -- 将 goctl 添加到需要安装的解析器列表中
      vim.list_extend(opts.ensure_installed, { "goctl" })

      -- 定义自定义解析器
      local parser_config = require("nvim-treesitter.parsers").get_parser_configs()
      parser_config.goctl = {
        install_info = {
          url = "https://github.com/chaozwn/tree-sitter-goctl",
          files = { "src/parser.c" },
          -- 如果需要，可以指定分支等其他字段
          -- branch = "main",
        },
        filetype = "goctl", -- 这是将解析器与文件类型关联的关键
      }
    end,
  },

  -- 在 lua/plugins/go-zero.lua 文件中，扩展 return 表
  --... 上述 nvim-treesitter 配置...
  {
    "neovim/neovim", -- 这是一个虚拟的插件项，仅用于执行配置
    config = function()
      -- 将.api 文件后缀映射到 goctl 文件类型
      vim.filetype.add({
        extension = {
          api = "goctl",
        },
      })
    end,
  },
}
