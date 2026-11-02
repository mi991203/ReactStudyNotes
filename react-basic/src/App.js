// 项目的根组件

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

export default App;
