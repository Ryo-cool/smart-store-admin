import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { deliveriesApi } from '@/lib/api/deliveries';

interface DeliveryStatusDialogProps {
  deliveryId: string;
  currentStatus: string;
  trigger: React.ReactNode;
}

export function DeliveryStatusDialog({
  deliveryId,
  currentStatus,
  trigger,
}: DeliveryStatusDialogProps) {
  const [open, setOpen] = useState(false);
  const [status, setStatus] = useState(currentStatus);
  const queryClient = useQueryClient();

  const { mutate, isPending } = useMutation({
    mutationFn: () => deliveriesApi.updateDeliveryStatus(deliveryId, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['delivery', deliveryId] });
      queryClient.invalidateQueries({ queryKey: ['delivery-history', deliveryId] });
      setOpen(false);
    },
  });

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{trigger}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>配送ステータスの更新</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <Select value={status} onValueChange={setStatus}>
            <SelectTrigger>
              <SelectValue placeholder="ステータスを選択" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="配送準備中">配送準備中</SelectItem>
              <SelectItem value="配送中">配送中</SelectItem>
              <SelectItem value="配送完了">配送完了</SelectItem>
              <SelectItem value="配送失敗">配送失敗</SelectItem>
            </SelectContent>
          </Select>
          <div className="flex justify-end gap-2">
            <Button variant="outline" onClick={() => setOpen(false)}>
              キャンセル
            </Button>
            <Button onClick={() => mutate()} disabled={isPending}>
              {isPending ? '更新中...' : '更新'}
            </Button>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
} 