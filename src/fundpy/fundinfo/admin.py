from django.contrib import admin
from fundinfo.models import *

# Register your models here.


class TransationAdmin(admin.ModelAdmin):
    list_display = ['id','code','date','cost','units']
    search_fields = ['code']
    ordering = ('code', '-date')

admin.site.register(Transation, TransationAdmin)