import { useCallback, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import { IconUpload, IconX } from '@tabler/icons-react';
import { cn } from '@/lib/utils';
import { Button } from './button';

interface ImageUploadProps {
  value?: string[];
  onChange?: (value: string[]) => void;
  maxFiles?: number;
  className?: string;
}

export function ImageUpload({
  value = [],
  onChange,
  maxFiles = 5,
  className,
}: ImageUploadProps) {
  const [images, setImages] = useState<string[]>(value);

  const onDrop = useCallback(
    (acceptedFiles: File[]) => {
      const newImages = acceptedFiles.map((file) => URL.createObjectURL(file));
      const updatedImages = [...images, ...newImages].slice(0, maxFiles);
      setImages(updatedImages);
      onChange?.(updatedImages);
    },
    [images, maxFiles, onChange]
  );

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      'image/*': ['.png', '.jpg', '.jpeg', '.gif'],
    },
    maxFiles,
  });

  const removeImage = (index: number) => {
    const updatedImages = images.filter((_, i) => i !== index);
    setImages(updatedImages);
    onChange?.(updatedImages);
  };

  return (
    <div className={className}>
      <div
        {...getRootProps()}
        className={cn(
          'relative flex h-32 cursor-pointer flex-col items-center justify-center rounded-lg border-2 border-dashed border-gray-300 bg-gray-50 hover:bg-gray-100',
          isDragActive && 'border-primary bg-primary/10',
          className
        )}
      >
        <input {...getInputProps()} />
        <IconUpload className="mb-2 h-6 w-6 text-gray-500" />
        <p className="text-sm text-gray-500">
          {isDragActive
            ? 'ドロップしてアップロード'
            : '画像をドラッグ＆ドロップまたはクリックしてアップロード'}
        </p>
        <p className="mt-1 text-xs text-gray-500">
          PNG, JPG, GIF（最大{maxFiles}枚）
        </p>
      </div>

      {images.length > 0 && (
        <div className="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5">
          {images.map((image, index) => (
            <div key={index} className="relative aspect-square">
              <img
                src={image}
                alt={`アップロード画像 ${index + 1}`}
                className="h-full w-full rounded-lg object-cover"
              />
              <Button
                variant="destructive"
                size="icon"
                className="absolute -right-2 -top-2 h-6 w-6"
                onClick={() => removeImage(index)}
              >
                <IconX className="h-4 w-4" />
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
} 