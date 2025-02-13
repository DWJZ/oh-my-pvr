"use client";

import React, { useState, ChangeEvent, FormEvent } from "react";

function App() {
  const [files, setFiles] = useState<File[] | null>(null);
  const [uploadDir, setUploadDir] = useState<string>("./uploads");
  const [message, setMessage] = useState<string[]>([]);

  const handleFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const fileList = Array.from(e.target.files);
      setFiles(fileList);
    }
  };

  const handleUploadDirChange = (e: ChangeEvent<HTMLInputElement>) => {
    setUploadDir(e.target.value);
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
    // Append the chosen upload directory
    formData.append("uploadDir", uploadDir);

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
    <div className="min-h-screen bg-gray-100 flex items-center justify-center p-4">
      <div className="bg-white shadow-md rounded-lg p-6 w-full max-w-md">
        <h1 className="text-3xl font-bold text-center">File Upload System</h1>
        <form onSubmit={handleSubmit} className="mt-6">
          <div className="mb-4">
            <label className="block text-gray-700 mb-1">
              Upload Directory
            </label>
            <input
              type="text"
              value={uploadDir}
              onChange={handleUploadDirChange}
              className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Enter upload directory"
            />
          </div>
          <div className="mb-4">
            <label className="block text-gray-700 mb-1">Select Files</label>
            <input
              type="file"
              multiple
              onChange={handleFileChange}
              className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <button
            type="submit"
            className="mt-4 w-full bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 transition duration-300"
          >
            Upload
          </button>
        </form>
        {message.length > 0 && (
          <div className="mt-4">
            {message.map((msg, index) => (
              <p key={index} className="text-red-500 text-sm">
                {msg}
              </p>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

export default App;