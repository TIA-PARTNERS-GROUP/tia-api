"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.app = void 0;
const express_1 = __importDefault(require("express"));
const swagger_ui_express_1 = __importDefault(require("swagger-ui-express"));
// @ts-ignore
const routes_1 = require("./routes");
exports.app = (0, express_1.default)();
const PORT = process.env.PORT || 8000;
exports.app.use(express_1.default.json());
(0, routes_1.RegisterRoutes)(exports.app);
try {
    const swaggerDocument = require('../docs/swagger.json');
    exports.app.use('/api-docs', swagger_ui_express_1.default.serve, swagger_ui_express_1.default.setup(swaggerDocument));
}
catch (error) {
    console.error('Error: Unable to load swagger.json. Please run the docs build command.');
}
exports.app.get('/health', (req, res) => {
    res.status(200).json({ status: 'OK', timestamp: new Date().toISOString() });
});
exports.app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
    console.log(`Swagger docs available at http://localhost:${PORT}/api-docs`);
});
