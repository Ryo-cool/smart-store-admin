import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useQueryClient } from '@tanstack/react-query';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
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
import { inventoryApi } from '@/lib/api/inventory';

const updateSchema = z.object({
  type: z.enum(['入庫', '出庫', '在庫調整']),
  quantity: z.number().min(1, '数量は1以上で入力してください'),
  reason: z.string().min(1, '理由を入力してください'),
  note: z.string().optional(),
});

type UpdateForm = z.infer<typeof updateSchema>;

interface InventoryUpdateDialogProps {
  productId: string;
  productName: string;
  trigger: React.ReactNode;
}

export function InventoryUpdateDialog({
  productId,
  productName,
  trigger,
}: InventoryUpdateDialogProps) {
  const [open, setOpen] = useState(false);
  const queryClient = useQueryClient();
  const { toast } = useToast();

  const form = useForm<UpdateForm>({
    resolver: zodResolver(updateSchema),
    defaultValues: {
      type: '入庫',
      quantity: 1,
    },
  });

  const onSubmit = async (data: UpdateForm) => {
    try {
      await inventoryApi.updateInventory({
        productId,
        ...data,
      });
      toast({
        title: '在庫を更新しました',
        description: '在庫の更新が完了しました。',
      });
      // 在庫履歴と商品情報を再取得
      queryClient.invalidateQueries({ queryKey: ['product', productId] });
      queryClient.invalidateQueries({ queryKey: ['inventory', productId] });
      setOpen(false);
      form.reset();
    } catch (error: unknown) {
      const errorMessage = error instanceof Error ? error.message : '在庫の更新に失敗しました。';
      toast({
        variant: 'destructive',
        title: 'エラー',
        description: errorMessage,
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>在庫更新</DialogTitle>
          <DialogDescription>{productName}の在庫を更新します。</DialogDescription>
        </DialogHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <FormField
              control={form.control}
              name="type"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>種別</FormLabel>
                  <Select onValueChange={field.onChange} defaultValue={field.value}>
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="種別を選択" />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      <SelectItem value="入庫">入庫</SelectItem>
                      <SelectItem value="出庫">出庫</SelectItem>
                      <SelectItem value="在庫調整">在庫調整</SelectItem>
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="quantity"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>数量</FormLabel>
                  <FormControl>
                    <Input
                      type="number"
                      min={1}
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
              name="reason"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>理由</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>
                  <FormDescription>
                    在庫を更新する理由を入力してください。
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="note"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>備考</FormLabel>
                  <FormControl>
                    <Textarea {...field} />
                  </FormControl>
                  <FormDescription>
                    必要に応じて補足情報を入力してください。
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex justify-end gap-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => setOpen(false)}
              >
                キャンセル
              </Button>
              <Button type="submit">更新</Button>
            </div>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
} 