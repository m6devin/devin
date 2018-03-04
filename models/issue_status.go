package models

import "time"

// IssueStatus contains:
//
// pending : معلق، هنوز کسی آن را پیگیری نکرده است. این حالت نمایانگر باز بودن مساله است
//
//  in progress : زمانیکه تیم شروع به بررسی و پیگیری آن کنند به در حال پیگیری تغییر حالت میدهد و دقت داشته باشید که کامنت هایی که در یک مساله ارسال میشود الزاما نشانگر در حال پیگیری بودن نمیباشد و این حالت باید توسط ادمین تعیین شود.
//
// closed : مساله بررسی و به نتیجه رسید که تنیجه نهایی آن در صورتیکه مساله غیر گیت باشد باید مشروح نوشته شود و سپس بسته شود. در مسايل گیت غالبا یا با آپدیت کد مشکل برطرف میشود که در این حالت مسئله به لیست باگ ها منتقل میشود و از کانالی دیگر پیگیری میشود و یا اینکه نیاز به راهنمایی جهت حل مشکل در بخشی است که در کامنت های افراد توضیحات یافت میشود.
type IssueStatus struct {
	ID        uint
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
