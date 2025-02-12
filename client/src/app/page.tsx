"use client";

import React, { useState, ChangeEvent, FormEvent } from "react";

function App() {
  const [files, setFiles] = useState<File[] | null>(null);
  const [message, setMessage] = useState<string[]>([]);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const fileList = Array.from(e.target.files);
      setFiles(fileList);
    }
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!files || files.length === 0) {
      setMessage(["Please select at least one file."]);
      return;
    }

    const formData = new FormData();
    files.forEach((file) => {
      formData.append("file", file);
    });

    try {
      const response = await fetch("http://localhost:8080/upload", {
        method: "POST",
        body: formData,
      });

      if (response.ok) {
        const result = await response.text();
        setMessage(result.split("\n"));
      } else {
        setMessage(["File upload failed."]);
      }
    } catch (error) {
      setMessage(["Error uploading file."]);
    }
  };

  return (
    <div className="App">
      <h1>File Upload System</h1>
      <form onSubmit={handleSubmit}>
        <input type="file" multiple onChange={handleFileChange} />
        <button type="submit">Upload</button>
      </form>
      {message.map((msg, index) => (
        <p key={index}>{msg}</p>
      ))}
    </div>
  );
}

export default App;
