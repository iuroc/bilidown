import js from '@eslint/js'
import globals from 'globals'
import tseslint from 'typescript-eslint'
import { defineConfig, globalIgnores } from 'eslint/config'
import stylistic from '@stylistic/eslint-plugin'
import importPlugin from 'eslint-plugin-import'

export default defineConfig([
    globalIgnores(['dist/', 'node_modules/']),
    {
        files: ['**/*.{js,mjs,cjs,ts,mts,cts}'],
        plugins: {
            js,
            '@stylistic': stylistic,
            'import': importPlugin
        },
        extends: ['js/recommended'],
        languageOptions: {
            globals: globals.browser
        },
        rules: {
            // 文件末尾不允许有空行
            '@stylistic/eol-last': ['error', 'never'],
            // 函数定义只能使用箭头函数
            'func-style': ['error', 'expression'],
            // 只有需要使用 await 时才使用 async
            'require-await': 'error',
            // async 函数必须使用 await 进行等待
            '@typescript-eslint/no-floating-promises': 'error',
            'no-restricted-imports': [
                'error',
                {
                    patterns: [
                        './*', // 禁止 ./ 开头的导入
                        '../*' // 禁止 ../ 开头的导入
                    ]
                }
            ]
        }
    },
    tseslint.configs.recommendedTypeChecked,
    {
        languageOptions: {
            parserOptions: {
                projectService: true
            }
        }
    },
    stylistic.configs.customize({
        // 代码缩进为 4 个空格
        indent: 4,
        // 不允许使用尾随逗号
        commaDangle: 'never',
        // 字符串使用单引号
        quotes: 'single',
        // 禁止行尾分号
        semi: false
    })
])

// Note: VSCode ESLint 插件版本为 3.0.20，插件标识符为 dbaeumer.vscode-eslint