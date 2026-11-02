# useState基础使用

useState 是一个 React HOOK 函数，允许想组件添加一个状态变量，从而控制组件的渲染效果。状态变量发生变化后，视图UI也会跟着变化（数据驱动试视图）

```js
const [count, setState] = useState(0)

1. useState是一个函数，返回值是一个数组
2. 数组第一个返回值是状态变量，第二个返回值是修改状态变量的函数
3. useState的参数将作为 count 的初始值
```

## 案例

```js
import { useState } from "react";

function App() {

  const [count, setCount] = useState(0)

  const handleClick = () => {
    setCount(count+1)
  }
  
  return (
    <div>
      <button onClick={handleClick}>{count}</button>
    </div>

  );
}
```

## 状态不可变

在 React 中，状态被认为是只读的，我们始终应该替换它而不是修改它，直接修改状态不能引发视图的更新