// 项目的根组件

// App -> index.js -> public/index.html(root)

function App() {

  // const handleClick = () => {
  //   console.log("按钮被点击了")
  // }

  // const handleClick = (e) => {
  //   console.log("按钮被点击了", e)
  // }

  const handleClick = (param, e) => {
    console.log("按钮被点击了", param, e)
  }
  
  return (
    <div>
      <button onClick={(e) => handleClick("Jack", e)}>按钮</button>
    </div>

  );
}

export default App;
