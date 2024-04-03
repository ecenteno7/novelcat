import logo from './logo.svg';
import './App.css';
function handleClick() {
  fetch('https://www.dev.mybooks.tech/seedDb')
    .then(response => response.json())
    .then(data => {
      // Handle the response data here
      console.log(data);
    })
    .catch(error => {
      // Handle any errors that occur during the request
      console.error(error);
    });
}
function App() {
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <h1>NOVEL CAT</h1>
        <p>
          Every book has a story.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Download the App
        </a>
        <button onClick={handleClick}>Seed Database</button>
      </header>
    </div>
  );
}

export default App;
