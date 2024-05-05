import logo from './logo.svg';
import './App.css';
import { seedDb } from './services/bookApi';
function handleClick() {
  seedDb();
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
        {/*<button onClick={handleClick}>Seed Database</button>*/}
      </header>
    </div>
  );
}

export default App;
