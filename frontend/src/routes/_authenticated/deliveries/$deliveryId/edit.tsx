import { Link, useParams, createFileRoute } from '@tanstack/react-router';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { IconArrowLeft } from '@tabler/icons-react';
import { useToast } from '@/components/ui/use-toast';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Skeleton } from '@/components/ui/skeleton';
import { deliveriesApi } from '@/lib/api/deliveries';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

export const Route = createFileRoute('/_authenticated/deliveries/$deliveryId/edit')({
  component: DeliveryEditPage,
});

const deliverySchema = z.object({
  deliveryType: z.string(),
  address: z.string().min(1, '配送先を入力してください'),
  estimatedDeliveryTime: z.string(),
  status: z.string(),
  notes: z.string().optional(),
});

type DeliveryFormValues = z.infer<typeof deliverySchema>;

function DeliveryEditPage() {
  const { deliveryId } = useParams({ from: '/_authenticated/deliveries/$deliveryId/edit' });
  const { toast } = useToast();
  const queryClient = useQueryClient();

  const { data: delivery, isLoading } = useQuery({
    queryKey: ['delivery', deliveryId],
    queryFn: () => deliveriesApi.getDelivery(deliveryId),
  });

  const form = useForm<DeliveryFormValues>({
    resolver: zodResolver(deliverySchema),
    values: delivery
      ? {
          deliveryType: delivery.deliveryType,
          address: delivery.address,
          estimatedDeliveryTime: delivery.estimatedDeliveryTime,
          status: delivery.status,
          notes: delivery.notes ?? '',
        }
      : undefined,
  });

  const { mutate: updateDelivery, isPending } = useMutation({
    mutationFn: (data: DeliveryFormValues) =>
      deliveriesApi.updateDelivery(deliveryId, data),
    onSuccess: () => {
      toast({
        title: '配送情報を更新しました',
        description: '配送情報が正常に更新されました。',
      });
      queryClient.invalidateQueries({ queryKey: ['delivery', deliveryId] });
    },
    onError: (error) => {
      toast({
        title: 'エラーが発生しました',
        description: error instanceof Error ? error.message : '不明なエラーです',
        variant: 'destructive',
      });
    },
  });

  if (isLoading) {
    return <Skeleton className="h-96 w-full" />;
  }

  if (!delivery) {
    return <div className="text-center">配送情報が見つかりませんでした</div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Link to="..">
          <Button variant="ghost" size="icon">
            <IconArrowLeft className="h-4 w-4" />
          </Button>
        </Link>
        <div>
          <h1 className="text-2xl font-bold">配送情報の編集</h1>
          <p className="text-sm text-gray-500">配送ID: {delivery.id}</p>
        </div>
      </div>

      <Form {...form}>
        <form
          onSubmit={form.handleSubmit((data) => updateDelivery(data))}
          className="space-y-6"
        >
          <FormField
            control={form.control}
            name="deliveryType"
            render={({ field }) => (
              <FormItem>
                <FormLabel>配送方法</FormLabel>
                <Select
                  value={field.value}
                  onValueChange={field.onChange}
                  disabled={isPending}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="配送方法を選択" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="ドローン">ドローン配送</SelectItem>
                    <SelectItem value="ロボット">ロボット配送</SelectItem>
                    <SelectItem value="通常配送">通常配送</SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="address"
            render={({ field }) => (
              <FormItem>
                <FormLabel>配送先</FormLabel>
                <FormControl>
                  <Input {...field} disabled={isPending} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="estimatedDeliveryTime"
            render={({ field }) => (
              <FormItem>
                <FormLabel>配送予定日時</FormLabel>
                <FormControl>
                  <Input
                    type="datetime-local"
                    {...field}
                    disabled={isPending}
                  />
                </FormControl>
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
                  value={field.value}
                  onValueChange={field.onChange}
                  disabled={isPending}
                >
                  <FormControl>
                    <SelectTrigger>
                      <SelectValue placeholder="ステータスを選択" />
                    </SelectTrigger>
                  </FormControl>
                  <SelectContent>
                    <SelectItem value="配送準備中">配送準備中</SelectItem>
                    <SelectItem value="配送中">配送中</SelectItem>
                    <SelectItem value="配送完了">配送完了</SelectItem>
                    <SelectItem value="配送失敗">配送失敗</SelectItem>
                  </SelectContent>
                </Select>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="notes"
            render={({ field }) => (
              <FormItem>
                <FormLabel>備考</FormLabel>
                <FormControl>
                  <Textarea {...field} disabled={isPending} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="flex justify-end gap-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => window.history.back()}
              disabled={isPending}
            >
              キャンセル
            </Button>
            <Button type="submit" disabled={isPending}>
              {isPending ? '更新中...' : '更新する'}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
} 