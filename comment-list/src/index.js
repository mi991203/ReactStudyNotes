
// React 必要的核心包
import React from 'react';
import ReactDOM from 'react-dom/client';

// 导入项目的根组件
import App from './App';

// 把App根组件渲染到id为root的dom上（public/index.html)
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);
