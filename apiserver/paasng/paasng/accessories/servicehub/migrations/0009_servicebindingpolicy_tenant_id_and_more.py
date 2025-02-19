# Generated by Django 4.2.16 on 2025-02-18 09:55

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('servicehub', '0008_remoteserviceengineappattachment_tenant_id_and_more'),
    ]

    operations = [
        migrations.AddField(
            model_name='servicebindingpolicy',
            name='tenant_id',
            field=models.CharField(default='default', help_text='本条数据的所属租户', max_length=32, verbose_name='租户 ID'),
        ),
        migrations.AddField(
            model_name='servicebindingprecedencepolicy',
            name='tenant_id',
            field=models.CharField(default='default', help_text='本条数据的所属租户', max_length=32, verbose_name='租户 ID'),
        ),
        migrations.AlterUniqueTogether(
            name='servicebindingpolicy',
            unique_together={('tenant_id', 'service_id')},
        ),
        migrations.AlterUniqueTogether(
            name='servicebindingprecedencepolicy',
            unique_together={('tenant_id', 'service_id', 'priority')},
        ),
    ]
