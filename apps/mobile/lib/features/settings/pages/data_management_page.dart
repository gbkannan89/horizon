import 'package:flutter/material.dart' hide ThemeMode; // Avoid potential theme conflicts
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';
import 'package:file_picker/file_picker.dart';
import 'package:dio/dio.dart';

class DataManagementPage extends ConsumerWidget {
  const DataManagementPage({super.key});

  Future<void> _export(WidgetRef ref, BuildContext context) async {
    try {
      final api = ref.read(apiClientProvider);
      await api.get('/data/export');
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Data exported successfully'), behavior: SnackBarBehavior.floating));
      }
    } catch (e) {
      if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Export failed: $e'), behavior: SnackBarBehavior.floating));
    }
  }

  Future<void> _backup(WidgetRef ref, BuildContext context) async {
    try {
      final api = ref.read(apiClientProvider);
      await api.post('/data/backup');
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Backup created'), behavior: SnackBarBehavior.floating));
      }
    } catch (e) {
      if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Backup failed: $e'), behavior: SnackBarBehavior.floating));
    }
  }

  void _showImportSheet(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (_) => const _ImportBottomSheet(),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        title: Text('Data Management', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: Container(
        decoration: BoxDecoration(
          gradient: AppColors.backgroundGradient(isDark),
        ),
        child: ListView(
          padding: const EdgeInsets.all(AppSpacing.md),
          children: [
            _card(theme, Icons.download_rounded, 'Export My Data', 'Download your financial data', 'Available', () => _export(ref, context)),
            _card(theme, Icons.upload_file_rounded, 'Import Data', 'Import bank statements (Excel, CSV, PDF)', 'Active', () => _showImportSheet(context, ref)),
            _card(theme, Icons.backup_rounded, 'Backup', 'Create a backup of your data', 'Available', () => _backup(ref, context)),
          ],
        ),
      ),
    );
  }

  Widget _card(ThemeData t, IconData icon, String title, String subtitle, String status, VoidCallback? onTap) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: GlassCard(
        onTap: onTap,
        child: Opacity(
          opacity: onTap != null ? 1.0 : 0.5,
          child: ListTile(
            contentPadding: EdgeInsets.zero,
            leading: Icon(icon, color: AppColors.teal500),
            title: Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
            subtitle: Text(subtitle),
            trailing: Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
              decoration: BoxDecoration(
                color: status == 'Active' ? AppColors.teal500.withOpacity(0.15) : AppColors.slate500.withOpacity(0.15),
                borderRadius: BorderRadius.circular(AppRadius.xs),
              ),
              child: Text(
                status,
                style: TextStyle(
                  color: status == 'Active' ? AppColors.teal500 : AppColors.slate600,
                  fontSize: 10,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class _ImportBottomSheet extends ConsumerStatefulWidget {
  const _ImportBottomSheet();

  @override
  ConsumerState<_ImportBottomSheet> createState() => _ImportBottomSheetState();
}

class _ImportBottomSheetState extends ConsumerState<_ImportBottomSheet> {
  String _selectedFormat = 'excel'; // 'excel', 'csv', 'pdf'
  PlatformFile? _selectedFile;
  bool _uploading = false;
  Map<String, dynamic>? _result;
  String? _error;

  Future<void> _pickFile() async {
    List<String> allowedExtensions;
    if (_selectedFormat == 'excel') {
      allowedExtensions = ['xlsx', 'xls'];
    } else if (_selectedFormat == 'csv') {
      allowedExtensions = ['csv'];
    } else {
      allowedExtensions = ['pdf'];
    }

    try {
      final result = await FilePicker.platform.pickFiles(
        type: FileType.custom,
        allowedExtensions: allowedExtensions,
      );
      if (result != null && result.files.isNotEmpty) {
        setState(() {
          _selectedFile = result.files.first;
          _error = null;
          _result = null;
        });
      }
    } catch (e) {
      setState(() => _error = 'Failed to pick file: $e');
    }
  }

  Future<void> _upload() async {
    if (_selectedFile == null) return;
    setState(() {
      _uploading = true;
      _error = null;
      _result = null;
    });

    try {
      final api = ref.read(apiClientProvider);
      
      MultipartFile multipartFile;
      if (_selectedFile!.path != null) {
        multipartFile = await MultipartFile.fromFile(
          _selectedFile!.path!,
          filename: _selectedFile!.name,
        );
      } else if (_selectedFile!.bytes != null) {
        multipartFile = MultipartFile.fromBytes(
          _selectedFile!.bytes!,
          filename: _selectedFile!.name,
        );
      } else {
        throw 'No file path or data available';
      }

      final formData = FormData.fromMap({
        'file': multipartFile,
      });

      final response = await api.post(
        '/import/$_selectedFormat',
        data: formData,
        options: Options(
          headers: {'Content-Type': 'multipart/form-data'},
        ),
      );

      setState(() {
        _result = response.data as Map<String, dynamic>;
        _uploading = false;
      });
    } catch (e) {
      setState(() {
        _error = 'Upload failed: $e';
        _uploading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Container(
      decoration: BoxDecoration(
        gradient: AppColors.backgroundGradient(isDark),
        borderRadius: const BorderRadius.vertical(top: Radius.circular(AppRadius.xl)),
      ),
      padding: EdgeInsets.only(
        top: 24,
        left: 24,
        right: 24,
        bottom: 24 + MediaQuery.of(context).viewInsets.bottom,
      ),
      child: SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('Import Bank Statement', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                IconButton(
                  icon: const Icon(Icons.close_rounded),
                  onPressed: () => Navigator.pop(context),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Text('Select Format', style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600, color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 8),
            Row(
              children: [
                _formatChip('excel', 'Excel (.xlsx)'),
                const SizedBox(width: 8),
                _formatChip('csv', 'CSV (.csv)'),
                const SizedBox(width: 8),
                _formatChip('pdf', 'PDF (.pdf)'),
              ],
            ),
            const SizedBox(height: 24),
            GestureDetector(
              onTap: _uploading ? null : _pickFile,
              child: GlassCard(
                padding: const EdgeInsets.symmetric(vertical: 32, horizontal: 16),
                child: Column(
                  children: [
                    Icon(
                      _selectedFile != null ? Icons.insert_drive_file_rounded : Icons.cloud_upload_rounded,
                      size: 48,
                      color: AppColors.teal500,
                    ),
                    const SizedBox(height: 16),
                    Text(
                      _selectedFile != null ? _selectedFile!.name : 'Tap to browse bank statement file',
                      textAlign: TextAlign.center,
                      style: const TextStyle(fontWeight: FontWeight.bold),
                    ),
                    if (_selectedFile != null) ...[
                      const SizedBox(height: 4),
                      Text(
                        'Size: ${(_selectedFile!.size / 1024).toStringAsFixed(1)} KB',
                        style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                      ),
                    ],
                  ],
                ),
              ),
            ),
            const SizedBox(height: 24),
            if (_error != null) ...[
              Text(
                _error!,
                style: TextStyle(color: AppColors.red500, fontWeight: FontWeight.bold),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
            ],
            if (_result != null) ...[
              GlassCard(
                child: Column(
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(Icons.check_circle_rounded, color: AppColors.teal500),
                        const SizedBox(width: 8),
                        Text('Import Completed', style: TextStyle(color: AppColors.teal500, fontWeight: FontWeight.bold)),
                      ],
                    ),
                    const SizedBox(height: 12),
                    Text(
                      'Processed ${_result!['total']} transactions.',
                      style: theme.textTheme.bodyMedium,
                    ),
                    Text(
                      'Created ${_result!['events_created']} events.',
                      style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 24),
            ],
            if (_uploading)
              const Center(
                child: Padding(
                  padding: EdgeInsets.all(16.0),
                  child: CircularProgressIndicator(color: AppColors.teal500),
                ),
              )
            else
              FilledButton(
                onPressed: _selectedFile == null ? null : _upload,
                style: FilledButton.styleFrom(
                  backgroundColor: AppColors.teal500,
                  minimumSize: const Size(double.infinity, 48),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppRadius.md)),
                ),
                child: const Text('Upload & Process Statement', style: TextStyle(fontWeight: FontWeight.bold)),
              ),
          ],
        ),
      ),
    );
  }

  Widget _formatChip(String value, String label) {
    final isSelected = _selectedFormat == value;
    return Expanded(
      child: ChoiceChip(
        label: Text(label),
        selected: isSelected,
        onSelected: _uploading ? null : (selected) {
          if (selected) {
            setState(() {
              _selectedFormat = value;
              _selectedFile = null;
              _result = null;
              _error = null;
            });
          }
        },
        selectedColor: AppColors.teal500.withOpacity(0.15),
        checkmarkColor: AppColors.teal500,
        labelStyle: TextStyle(
          color: isSelected ? AppColors.teal500 : null,
          fontWeight: isSelected ? FontWeight.bold : null,
          fontSize: 12,
        ),
      ),
    );
  }
}
