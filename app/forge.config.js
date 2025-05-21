"use strict";
var _a;
Object.defineProperty(exports, "__esModule", { value: true });
var plugin_fuses_1 = require("@electron-forge/plugin-fuses");
var fuses_1 = require("@electron/fuses");
exports.default = {
    packagerConfig: {
        asar: true,
        executableName: 'bilidown',
        icon: '../web/public/favicon'
    },
    rebuildConfig: {},
    makers: [
        // {
        //   name: '@electron-forge/maker-squirrel',
        //   config: {},
        // },
        {
            name: '@electron-forge/maker-zip',
            platforms: ['darwin', 'win32'],
        },
        {
            name: '@electron-forge/maker-deb',
            config: {},
        },
        {
            name: '@electron-forge/maker-rpm',
            config: {},
        },
    ],
    publishers: [
        {
            name: '@electron-forge/publisher-github',
            config: {
                repository: {
                    owner: 'iuroc',
                    name: 'bilidown'
                },
                draft: true
            }
        }
    ],
    plugins: [
        {
            name: '@electron-forge/plugin-auto-unpack-natives',
            config: {},
        },
        // Fuses are used to enable/disable various Electron functionality
        // at package time, before code signing the application
        new plugin_fuses_1.FusesPlugin((_a = {
                version: fuses_1.FuseVersion.V1
            },
            _a[fuses_1.FuseV1Options.RunAsNode] = false,
            _a[fuses_1.FuseV1Options.EnableCookieEncryption] = true,
            _a[fuses_1.FuseV1Options.EnableNodeOptionsEnvironmentVariable] = false,
            _a[fuses_1.FuseV1Options.EnableNodeCliInspectArguments] = false,
            _a[fuses_1.FuseV1Options.EnableEmbeddedAsarIntegrityValidation] = true,
            _a[fuses_1.FuseV1Options.OnlyLoadAppFromAsar] = true,
            _a)),
    ],
};
