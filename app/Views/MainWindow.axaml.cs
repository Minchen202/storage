using System;
using System.IO;
using System.Linq;
using System.Threading.Tasks;
using System.Timers;
using Avalonia.Controls;
using Avalonia.Interactivity;
using Avalonia.Media;
using Avalonia.Platform.Storage;

namespace P2PStorageApp;

public partial class MainWindow : Window
{
    private readonly ApiService _api;
    private Timer _heartbeatTimer;

    public MainWindow()
    {
        InitializeComponent();
        _api = new ApiService();
    }

    private async void OnLoginClick(object sender, RoutedEventArgs e)
    {
        var email = EmailBox.Text;
        var password = PasswordBox.Text;

        if (string.IsNullOrWhiteSpace(email) || string.IsNullOrWhiteSpace(password))
        {
            ShowAuthError("Please enter email and password");
            return;
        }

        var (success, token, error) = await _api.LoginAsync(email, password);
        
        if (success)
        {
            AppConfig.JwtToken = token;
            AppConfig.UserEmail = email;
            ShowMainApp();
            StartHeartbeat();
            await LoadUserData();
        }
        else
        {
            ShowAuthError(error ?? "Login failed");
        }
    }

    private async void OnRegisterClick(object sender, RoutedEventArgs e)
    {
        var email = EmailBox.Text;
        var password = PasswordBox.Text;

        if (string.IsNullOrWhiteSpace(email) || string.IsNullOrWhiteSpace(password))
        {
            ShowAuthError("Please enter email and password");
            return;
        }

        if (password.Length < 8)
        {
            ShowAuthError("Password must be at least 8 characters");
            return;
        }

        var (success, error) = await _api.RegisterAsync(email, password);
        
        if (success)
        {
            ShowAuthError("Registration successful! Please login.");
            AuthError.Foreground = new SolidColorBrush(Color.Parse("#10B981"));
        }
        else
        {
            ShowAuthError(error ?? "Registration failed");
        }
    }

    private void ShowAuthError(string message)
    {
        AuthError.Text = message;
        AuthError.IsVisible = true;
        AuthError.Foreground = new SolidColorBrush(Color.Parse("#EF4444"));
    }

    private void ShowMainApp()
    {
        LoginScreen.IsVisible = false;
        AppScreen.IsVisible = true;
        UserEmailText.Text = AppConfig.UserEmail;
    }

    private void StartHeartbeat()
    {
        _heartbeatTimer = new Timer(600000); // 10 minutes
        _heartbeatTimer.Elapsed += async (s, e) => await _api.SendHeartbeatAsync();
        _heartbeatTimer.Start();
        
        // Send first heartbeat immediately
        Task.Run(async () => await _api.SendHeartbeatAsync());
    }

    private async Task LoadUserData()
    {
        var stats = await _api.GetUserStatsAsync();
        if (stats != null)
        {
            StorageContributed.Text = $"{(long)stats.storage_contributed / 1073741824} GB";
            StorageQuota.Text = $"{(long)stats.storage_quota / 1073741824} GB";
            TotalUptime.Text = $"{(long)stats.total_uptime_seconds / 3600} hours";
            
            var used = (long)stats.storage_contributed / 1073741824.0;
            var quota = (long)stats.storage_quota / 1073741824.0;
            StorageUsageText.Text = $"{used:F2} GB / {quota:F2} GB";
            StorageUsageBar.Width = (used / quota) * 500;
        }

        await LoadFiles();
    }

    private async Task LoadFiles()
    {
        var files = await _api.GetFilesAsync();
        FilesList.Children.Clear();

        if (files != null)
        {
            foreach (var file in files)
            {
                var fileItem = CreateFileItem(
                    file.filename.ToString(),
                    $"{(long)file.file_size / 1048576} MB",
                    file.chunk_count.ToString(),
                    file.created_at.ToString()
                );
                FilesList.Children.Add(fileItem);
            }
        }
    }

    private Border CreateFileItem(string name, string size, string chunks, string uploaded)
    {
        var border = new Border
        {
            Background = new SolidColorBrush(Color.Parse("#1F2937")),
            BorderBrush = new SolidColorBrush(Color.Parse("#374151")),
            BorderThickness = new Avalonia.Thickness(1),
            CornerRadius = new Avalonia.CornerRadius(8),
            Padding = new Avalonia.Thickness(16),
            Margin = new Avalonia.Thickness(0, 0, 0, 8)
        };

        var grid = new Grid();
        grid.ColumnDefinitions.Add(new ColumnDefinition(GridLength.Star));
        grid.ColumnDefinitions.Add(new ColumnDefinition(GridLength.Auto));

        var infoStack = new StackPanel();
        
        var nameText = new TextBlock
        {
            Text = name,
            Foreground = Brushes.White,
            FontWeight = FontWeight.SemiBold,
            Margin = new Avalonia.Thickness(0, 0, 0, 4)
        };
        
        var detailsText = new TextBlock
        {
            Text = $"{size} · {chunks} chunks · Uploaded {uploaded}",
            Foreground = new SolidColorBrush(Color.Parse("#9CA3AF")),
            FontSize = 12
        };

        infoStack.Children.Add(nameText);
        infoStack.Children.Add(detailsText);

        var buttonStack = new StackPanel
        {
            Orientation = Avalonia.Layout.Orientation.Horizontal,
            Spacing = 8
        };

        var statusBadge = new Border
        {
            Background = new SolidColorBrush(Color.Parse("#10B98140")),
            CornerRadius = new Avalonia.CornerRadius(12),
            Padding = new Avalonia.Thickness(12, 6),
            Child = new TextBlock
            {
                Text = "distributed",
                Foreground = new SolidColorBrush(Color.Parse("#10B981")),
                FontSize = 11,
                FontWeight = FontWeight.Bold
            }
        };

        var downloadBtn = new Button
        {
            Content = "⬇",
            Padding = new Avalonia.Thickness(8),
            Background = new SolidColorBrush(Color.Parse("#374151")),
            Foreground = new SolidColorBrush(Color.Parse("#9CA3AF"))
        };

        var deleteBtn = new Button
        {
            Content = "🗑",
            Padding = new Avalonia.Thickness(8),
            Background = new SolidColorBrush(Color.Parse("#374151")),
            Foreground = new SolidColorBrush(Color.Parse("#EF4444"))
        };

        buttonStack.Children.Add(statusBadge);
        buttonStack.Children.Add(downloadBtn);
        buttonStack.Children.Add(deleteBtn);

        Grid.SetColumn(infoStack, 0);
        Grid.SetColumn(buttonStack, 1);

        grid.Children.Add(infoStack);
        grid.Children.Add(buttonStack);
        border.Child = grid;

        return border;
    }

    private void OnFilesTabClick(object sender, RoutedEventArgs e)
    {
        SwitchTab("files");
    }

    private void OnUploadTabClick(object sender, RoutedEventArgs e)
    {
        SwitchTab("upload");
    }

    private void OnStatsTabClick(object sender, RoutedEventArgs e)
    {
        SwitchTab("stats");
    }

    private void OnSettingsTabClick(object sender, RoutedEventArgs e)
    {
        SwitchTab("settings");
    }

    private void SwitchTab(string tab)
    {
        // Reset all tab buttons
        FilesTab.Background = new SolidColorBrush(Color.Parse("#374151"));
        FilesTab.Foreground = new SolidColorBrush(Color.Parse("#9CA3AF"));
        UploadTab.Background = new SolidColorBrush(Color.Parse("#374151"));
        UploadTab.Foreground = new SolidColorBrush(Color.Parse("#9CA3AF"));
        StatsTab.Background = new SolidColorBrush(Color.Parse("#374151"));
        StatsTab.Foreground = new SolidColorBrush(Color.Parse("#9CA3AF"));
        SettingsTab.Background = new SolidColorBrush(Color.Parse("#374151"));
        SettingsTab.Foreground = new SolidColorBrush(Color.Parse("#9CA3AF"));

        // Hide all content
        FilesContent.IsVisible = false;
        UploadContent.IsVisible = false;
        StatsContent.IsVisible = false;
        SettingsContent.IsVisible = false;

        // Show selected tab
        switch (tab)
        {
            case "files":
                FilesTab.Background = new SolidColorBrush(Color.Parse("#2563EB"));
                FilesTab.Foreground = Brushes.White;
                FilesContent.IsVisible = true;
                break;
            case "upload":
                UploadTab.Background = new SolidColorBrush(Color.Parse("#2563EB"));
                UploadTab.Foreground = Brushes.White;
                UploadContent.IsVisible = true;
                break;
            case "stats":
                StatsTab.Background = new SolidColorBrush(Color.Parse("#2563EB"));
                StatsTab.Foreground = Brushes.White;
                StatsContent.IsVisible = true;
                break;
            case "settings":
                SettingsTab.Background = new SolidColorBrush(Color.Parse("#2563EB"));
                SettingsTab.Foreground = Brushes.White;
                SettingsContent.IsVisible = true;
                break;
        }
    }

    private void OnUploadFileClick(object sender, RoutedEventArgs e)
    {
        SwitchTab("upload");
    }

    private async void OnSelectFileClick(object sender, RoutedEventArgs e)
    {
        var files = await StorageProvider.OpenFilePickerAsync(new FilePickerOpenOptions
        {
            Title = "Select File to Upload",
            AllowMultiple = false
        });

        if (files.Any())
        {
            var file = files[0];
            await UploadFile(file);
        }
    }

    private async Task UploadFile(IStorageFile file)
    {
        try
        {
            using var stream = await file.OpenReadAsync();
            var fileData = new byte[stream.Length];
            await stream.ReadAsync(fileData, 0, (int)stream.Length);

            // Simple chunking (4MB chunks)
            const int chunkSize = 4 * 1024 * 1024;
            var chunks = (int)Math.Ceiling((double)fileData.Length / chunkSize);

            for (int i = 0; i < chunks; i++)
            {
                var offset = i * chunkSize;
                var length = Math.Min(chunkSize, fileData.Length - offset);
                var chunkData = new byte[length];
                Array.Copy(fileData, offset, chunkData, 0, length);

                var hash = BitConverter.ToString(
                    System.Security.Cryptography.SHA256.HashData(chunkData)
                ).Replace("-", "").ToLower();

                var (success, error) = await _api.UploadChunkAsync(hash, chunkData);
                
                if (!success)
                {
                    ShowAuthError($"Upload failed: {error}");
                    return;
                }
            }

            await LoadFiles();
            SwitchTab("files");
        }
        catch (Exception ex)
        {
            ShowAuthError($"Upload error: {ex.Message}");
        }
    }

    private void OnSaveSettingsClick(object sender, RoutedEventArgs e)
    {
        AppConfig.ApiBaseUrl = ApiEndpointBox.Text;
        ShowAuthError("Settings saved!");
        AuthError.Foreground = new SolidColorBrush(Color.Parse("#10B981"));
    }
}