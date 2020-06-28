import React, { useState, useEffect } from "react";

import ItemList from "./components/itemList/ItemList";
import Controls from "./components/controls/controls";
import axios from "axios";

import "./App.css";

function App() {
  var [items, setItems] = useState([]);
  const [newItem, setNewItem] = useState("");

  useEffect(() => {
    getData();
  }, []);

  const getData = async () => {
    const res = await axios.get(`http://localhost:8080/items`);
    if (res.status !== 200) {
      throw new Error(
        `Failed to retrieve data from the server. 
        (Response Status = ${res.status} - ${res.statusText})`
      );
    }

    setItems(res.data);
  };

  const addNewItem = async () => {
    const item = { title: newItem };
    const res = await axios.post(`http://localhost:8080/items`, item);
    if (res.status !== 200) {
      throw new Error(
        `Failed to add new Item. 
        (Response Status = ${res.status} - ${res.statusText})`
      );
    }

    setItems(res.data);
    setNewItem("");
  };

  const clearAll = async () => {
    const res = await axios.delete(`http://localhost:8080/clearAll`);
    if (res.status !== 200) {
      throw new Error(
        `Failed to clear all items. 
        (Response Status = ${res.status} - ${res.statusText})`
      );
    }

    setItems(res.data);
  };

  return (
    <div className="App">
      <ItemList items={items} />
      <Controls
        item={newItem}
        allItems={items}
        onChange={(newItem: string) => setNewItem(newItem)}
        onAdd={addNewItem}
        onClear={clearAll}
      />
    </div>
  );
}

export default App;
