import React, { Component } from 'react';
import './App.css';
import FamilyMembers from './FamilyMembers';

class App extends Component {
  render() {
    return (
      <div className="App">
        <header className="App-header">
          <FamilyMembers />
        </header>
      </div>
    );
  }
}

export default App;
