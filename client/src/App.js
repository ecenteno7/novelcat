import logo from './logo.svg';
import './App.css';
import { seedDb } from './services/bookApi';
import TextField from '@mui/material/TextField';

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
        <TextField id="standard-basic" label="Standard" variant="standard" />
      </header>
    </div >
  );
}

export default App;
