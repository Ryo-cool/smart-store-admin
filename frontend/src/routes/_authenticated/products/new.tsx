import { Link, useNavigate } from '@tanstack/react-router';
import { IconArrowLeft } from '@tabler/icons-react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Button } from '@/components/ui/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useToast } from '@/hooks/use-toast';
import { productsApi } from '@/lib/api/products';
import { ImageUpload } from '@/components/ui/image-upload';

const productSchema = z.object({
  name: z.string().min(1, '商品名を入力してください'),
  sku: z.string().min(1, 'SKUを入力してください'),
  price: z.number().min(0, '価格は0以上で入力してください'),
  stock: z.number().min(0, '在庫数は0以上で入力してください'),
  description: z.string().optional(),
  category: z.string().min(1, 'カテゴリーを選択してください'),
  weight: z.string().optional(),
  dimensions: z.string().optional(),
  status: z.enum(['販売中', '在庫切れ', '入荷待ち', '在庫少']),
  images: z.array(z.string()).optional(),
});

type ProductForm = z.infer<typeof productSchema>;

const categories = [
  { value: 'coffee', label: 'コーヒー豆' },
  { value: 'tea', label: '茶葉' },
  { value: 'equipment', label: '器具・設備' },
  { value: 'accessories', label: 'アクセサリー' },
];

export default function NewProductPage() {
  const navigate = useNavigate();
  const { toast } = useToast();
  const form = useForm<ProductForm>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      status: '販売中',
    },
  });

  const onSubmit = async (data: ProductForm) => {
    try {
      await productsApi.createProduct(data);
      toast({
        title: '商品を登録しました',
        description: '商品の登録が完了しました。',
      });
      navigate({ to: '/products' });
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : '商品の登録に失敗しました。';
      toast({
        variant: 'destructive',
        title: 'エラー',
        description: errorMessage,
      });
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="/products">
          <Button variant="ghost" size="icon">
            <IconArrowLeft className="h-4 w-4" />
          </Button>
        </Link>
        <div>
          <h1 className="text-2xl font-bold">新規商品登録</h1>
          <p className="text-sm text-gray-500">新しい商品を登録</p>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <div className="grid gap-6 md:grid-cols-2">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>商品名</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="sku"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>SKU</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="price"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>価格</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      {...field}
                      onChange={(e) => field.onChange(Number(e.target.value))}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="stock"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>在庫数</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      {...field}
                      onChange={(e) => field.onChange(Number(e.target.value))}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="category"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>カテゴリー</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="カテゴリーを選択" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {categories.map((category) => (
                        <SelectItem key={category.value} value={category.value}>
                          {category.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="status"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ステータス</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    defaultValue={field.value}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="ステータスを選択" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="販売中">販売中</SelectItem>
                      <SelectItem value="在庫切れ">在庫切れ</SelectItem>
                      <SelectItem value="入荷待ち">入荷待ち</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          <FormField
            control={form.control}
            name="description"
            render={({ field }) => (
              <FormItem>
                <FormLabel>商品説明</FormLabel>
                <FormControl>
                  <Textarea {...field} />
                </FormControl>
                <FormDescription>商品の詳細な説明を入力してください。</FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="grid gap-6 md:grid-cols-2">
            <FormField
              control={form.control}
              name="weight"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>重量</FormLabel>
                  <FormControl>
                    <Input {...field} placeholder="例: 200g" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="dimensions"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>サイズ</FormLabel>
                  <FormControl>
                    <Input {...field} placeholder="例: 15 x 8 x 8 cm" />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>

          <FormField
            control={form.control}
            name="images"
            render={({ field }) => (
              <FormItem>
                <FormLabel>商品画像</FormLabel>
                <FormControl>
                  <ImageUpload
                    value={field.value}
                    onChange={field.onChange}
                    maxFiles={5}
                  />
                </FormControl>
                <FormDescription>
                  商品の画像を最大5枚までアップロードできます。
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex justify-end gap-4">
            <Link to="/products">
              <Button variant="outline">キャンセル</Button>
            </Link>
            <Button type="submit">登録</Button>
          </div>
        </form>
      </Form>
    </div>
  );
} 