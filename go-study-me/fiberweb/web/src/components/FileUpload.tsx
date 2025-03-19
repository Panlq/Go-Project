import React, { useState, useRef } from 'react';
import { Upload, Button, Progress, message, Table } from 'antd';
import { UploadOutlined } from '@ant-design/icons';
import SparkMD5 from 'spark-md5';
import axios from 'axios';

interface ChunkInfo {
  chunk: Blob;
  hash: string;
  filename: string;
  index: number;
  size: number;
  progress: number;
}

const CHUNK_SIZE = 2 * 1024 * 1024; // 2MB per chunk

const MAX_CONCURRENT_UPLOADS = 10; // 最大并发上传数
const MAX_RETRIES = 3; // 最大重试次数

const FileUpload: React.FC = () => {
  const [uploading, setUploading] = useState(false);
  const [fileChunks, setFileChunks] = useState<ChunkInfo[]>([]);
  const [uploadProgress, setUploadProgress] = useState(0);
  const abortController = useRef<AbortController | null>(null);
  const [retryCount, setRetryCount] = useState<Record<number, number>>({});

  const calculateHash = async (file: File): Promise<string> => {
    const spark = new SparkMD5.ArrayBuffer();
    const reader = new FileReader();
    const chunks = Math.ceil(file.size / CHUNK_SIZE);
    let currentChunk = 0;
    console.log(`[INFO] Starting hash calculation for file: ${file.name}, size: ${file.size}, total chunks: ${chunks}`);

    return new Promise((resolve, reject) => {
      const loadNext = () => {
        const start = currentChunk * CHUNK_SIZE;
        const end = Math.min(start + CHUNK_SIZE, file.size);
        const chunk = file.slice(start, end);
        console.log(`[INFO] Reading chunk ${currentChunk}/${chunks}, start: ${start}, end: ${end}, size: ${chunk.size}`);
        reader.readAsArrayBuffer(chunk);
      };

      reader.onload = (e) => {
        const chunkData = e.target?.result as ArrayBuffer;
        const chunkHash = new SparkMD5.ArrayBuffer().append(chunkData).end();
        console.log(`[INFO] Processing chunk ${currentChunk}/${chunks}, size: ${chunkData.byteLength}, chunk hash: ${chunkHash}`);
        spark.append(chunkData);
        currentChunk++;
        if (currentChunk < chunks) {
          loadNext();
        } else {
          const finalHash = spark.end();
          console.log(`[INFO] Final hash calculated: ${finalHash}, total chunks processed: ${chunks}`);
          resolve(finalHash);
        }
      };

      reader.onerror = (error) => {
        console.error(`[ERROR] Failed to read chunk ${currentChunk}:`, error);
        reject(reader.error);
      };
      loadNext();
    });
  };

  const createFileChunks = async (file: File) => {
    const chunks: ChunkInfo[] = [];
    const hash = await calculateHash(file);
    const chunksCount = Math.ceil(file.size / CHUNK_SIZE);

    for (let i = 0; i < chunksCount; i++) {
      const start = i * CHUNK_SIZE;
      const end = Math.min(start + CHUNK_SIZE, file.size);
      chunks.push({
        chunk: file.slice(start, end),
        hash: `${hash}`,
        filename: file.name,
        index: i,
        size: end - start,
        progress: 0
      });
    }

    setFileChunks(chunks);
    return { fileHash: hash, chunks };
  };

  const uploadChunk = async (chunk: ChunkInfo) => {
    const formData = new FormData();
    formData.append('chunk', chunk.chunk);
    formData.append('hash', chunk.hash);
    formData.append('filename', chunk.filename);
    formData.append('index', chunk.index.toString());
    formData.append('size', chunk.size.toString());

    try {
      const response = await axios.post('/api/upload/chunk', formData, {
        signal: abortController.current?.signal,
        timeout: 30000, // 30秒超时
        onUploadProgress: (progressEvent) => {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total!);
          setFileChunks(prev => prev.map(item => {
            if (item.hash === chunk.hash && item.index === chunk.index) {
              return { ...item, progress: percentCompleted };
            }
            return item;
          }));

          const totalProgress = fileChunks.reduce((acc, curr) => acc + curr.progress, 0) / fileChunks.length;
          setUploadProgress(Math.round(totalProgress));
        }
      });
      setRetryCount(prev => ({ ...prev, [chunk.index]: 0 }));
      return true;
    } catch (error: any) {
      if (axios.isCancel(error)) {
        console.log('Upload canceled');
      } else {
        const isEPIPEError = error.code === 'EPIPE' || error.message?.includes('EPIPE');
        const currentRetries = retryCount[chunk.index] || 0;
        
        if (currentRetries < MAX_RETRIES) {
          setRetryCount(prev => ({ ...prev, [chunk.index]: currentRetries + 1 }));
          const retryDelay = isEPIPEError ? 2000 : 1000 * (currentRetries + 1); // EPIPE错误增加等待时间
          console.log(`Retrying chunk ${chunk.index}, attempt ${currentRetries + 1}/${MAX_RETRIES}${isEPIPEError ? ' (EPIPE error)' : ''}`);
          await new Promise(resolve => setTimeout(resolve, retryDelay));
          return uploadChunk(chunk);
        } else {
          const errorMessage = error.response?.data?.error || 
            `Chunk ${chunk.index} upload failed after ${MAX_RETRIES} retries${isEPIPEError ? ' (EPIPE error)' : ''}`;
          message.error(errorMessage);
          console.error('Upload error:', error);
        }
      }
      return false;
    }
  };

  const mergeRequest = async (filename: string, hash: string, totalSize: number) => {
    try {
      await axios.post('/api/upload/merge', {
        filename,
        hash,
        size: totalSize
      });
      message.success('File uploaded successfully');
    } catch (error: any) {
      const errorMessage = error.response?.data?.error || 'File merge failed';
      message.error(errorMessage);
      console.error('Merge error:', errorMessage);
    }
  };

  const handleUpload = async (file: File) => {
    try {
      setUploading(true);
      abortController.current = new AbortController();

      const { fileHash, chunks } = await createFileChunks(file);

      // Get uploaded chunks and check if file is complete
      const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
      const checkResponse = await axios.get(`/api/upload/check?hash=${fileHash}&filename=${file.name}&totalChunks=${totalChunks}`);
      if (checkResponse.data.exists) {
        message.success('File already exists');
        return;
      }
      const uploadedChunks = new Set(checkResponse.data.uploadedChunks || []);

      // Update progress for existing chunks
      setFileChunks(prev => prev.map(chunk => ({
        ...chunk,
        progress: uploadedChunks.has(chunk.index) ? 100 : 0
      })));

      // Upload only chunks that haven't been uploaded
      const remainingChunks = chunks.filter(chunk => !uploadedChunks.has(chunk.index));
      if (remainingChunks.length === 0) {
        await mergeRequest(file.name, fileHash, file.size);
        return;
      }

      // 分批上传文件块
      const results: boolean[] = [];
      for (let i = 0; i < remainingChunks.length; i += MAX_CONCURRENT_UPLOADS) {
        const batch = remainingChunks.slice(i, i + MAX_CONCURRENT_UPLOADS);
        const batchResults = await Promise.all(batch.map(chunk => uploadChunk(chunk)));
        results.push(...batchResults);
      }

      if (results.every(Boolean)) {
        await mergeRequest(file.name, fileHash, file.size);
      }
    } catch (error) {
      message.error('Upload failed');
    } finally {
      setUploading(false);
      setUploadProgress(0);
      setFileChunks([]);
      abortController.current = null;
    }
  };

  const cancelUpload = () => {
    abortController.current?.abort();
    setUploading(false);
    setUploadProgress(0);
    setFileChunks([]);
  };

  const columns = [
    {
      title: '分片序号',
      dataIndex: 'index',
      key: 'index',
    },
    {
      title: '分片大小',
      dataIndex: 'size',
      key: 'size',
      render: (size: number) => `${(size / 1024).toFixed(2)} KB`,
    },
    {
      title: '上传进度',
      dataIndex: 'progress',
      key: 'progress',
      render: (progress: number) => <Progress percent={progress} />,
    },
  ];

  return (
    <div style={{ padding: 24 }}>
      <Upload
        customRequest={({ file }) => handleUpload(file as File)}
        showUploadList={false}
      >
        <Button icon={<UploadOutlined />} loading={uploading}>
          Select File
        </Button>
      </Upload>
      {uploading && (
        <>
          <Button onClick={cancelUpload} style={{ marginLeft: 16 }}>
            Cancel
          </Button>
          <Progress percent={uploadProgress} style={{ marginTop: 16 }} />
          <Table
            dataSource={fileChunks}
            columns={columns}
            rowKey="index"
            pagination={false}
            style={{ marginTop: 16 }}
          />
        </>
      )}
    </div>
  );
};

export default FileUpload;