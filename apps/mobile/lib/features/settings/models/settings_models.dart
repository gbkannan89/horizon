class UserProfile {
  final String id; final String email; final String? name; final String? phone;
  final String? dateOfBirth; final String country; final String baseCurrency;
  final String? timezone; final String? language; final String createdAt;

  UserProfile({required this.id, required this.email, this.name, this.phone,
    this.dateOfBirth, this.country = 'IN', this.baseCurrency = 'INR',
    this.timezone, this.language, required this.createdAt});

  factory UserProfile.fromJson(Map<String, dynamic> json) => UserProfile(
    id: json['id'] as String? ?? json['user_id'] as String? ?? '',
    email: json['email'] as String? ?? '',
    name: json['name'] as String?, phone: json['phone'] as String?,
    dateOfBirth: json['date_of_birth'] as String?,
    country: json['country'] as String? ?? 'IN',
    baseCurrency: json['base_currency'] as String? ?? 'INR',
    timezone: json['timezone'] as String?, language: json['language'] as String?,
    createdAt: json['created_at'] as String? ?? '',
  );

  Map<String, dynamic> toJson() => {
    if (name != null) 'name': name, if (phone != null) 'phone': phone,
    if (dateOfBirth != null) 'date_of_birth': dateOfBirth,
    if (timezone != null) 'timezone': timezone, if (language != null) 'language': language,
  };
}

class UserPreferences {
  String baseCurrency; String language; String theme; bool emailNotifications;
  bool pushNotifications; bool smsNotifications; bool digestEnabled;
  String digestFrequency; int dashboardDefaultView; bool timelineCompact;

  UserPreferences({
    this.baseCurrency = 'INR', this.language = 'en', this.theme = 'system',
    this.emailNotifications = true, this.pushNotifications = true,
    this.smsNotifications = false, this.digestEnabled = false,
    this.digestFrequency = 'weekly', this.dashboardDefaultView = 0,
    this.timelineCompact = false,
  });

  factory UserPreferences.fromJson(Map<String, dynamic> json) => UserPreferences(
    baseCurrency: json['base_currency'] as String? ?? 'INR',
    language: json['language'] as String? ?? 'en',
    theme: json['theme'] as String? ?? 'system',
    emailNotifications: json['email_notifications'] as bool? ?? true,
    pushNotifications: json['push_notifications'] as bool? ?? true,
    smsNotifications: json['sms_notifications'] as bool? ?? false,
    digestEnabled: json['digest_enabled'] as bool? ?? false,
    digestFrequency: json['digest_frequency'] as String? ?? 'weekly',
    dashboardDefaultView: (json['dashboard_default_view'] as num?)?.toInt() ?? 0,
    timelineCompact: json['timeline_compact'] as bool? ?? false,
  );

  Map<String, dynamic> toJson() => {
    'base_currency': baseCurrency, 'language': language, 'theme': theme,
    'email_notifications': emailNotifications, 'push_notifications': pushNotifications,
    'sms_notifications': smsNotifications, 'digest_enabled': digestEnabled,
    'digest_frequency': digestFrequency, 'dashboard_default_view': dashboardDefaultView,
    'timeline_compact': timelineCompact,
  };
}

class PrivacySettings {
  bool dataSharingEnabled; bool analyticsEnabled; bool personalizeEnabled;
  bool thirdPartySharing; bool marketingEnabled;

  PrivacySettings({
    this.dataSharingEnabled = true, this.analyticsEnabled = true,
    this.personalizeEnabled = true, this.thirdPartySharing = false,
    this.marketingEnabled = false,
  });

  factory PrivacySettings.fromJson(Map<String, dynamic> json) => PrivacySettings(
    dataSharingEnabled: json['data_sharing_enabled'] as bool? ?? true,
    analyticsEnabled: json['analytics_enabled'] as bool? ?? true,
    personalizeEnabled: json['personalize_enabled'] as bool? ?? true,
    thirdPartySharing: json['third_party_sharing'] as bool? ?? false,
    marketingEnabled: json['marketing_enabled'] as bool? ?? false,
  );

  Map<String, dynamic> toJson() => {
    'data_sharing_enabled': dataSharingEnabled, 'analytics_enabled': analyticsEnabled,
    'personalize_enabled': personalizeEnabled, 'third_party_sharing': thirdPartySharing,
    'marketing_enabled': marketingEnabled,
  };
}
