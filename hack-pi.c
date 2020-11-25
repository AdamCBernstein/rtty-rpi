/*
 * RTTY signal generator using GPIO interface of Raspberry Pi.
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <signal.h>
#include <ctype.h>
#include <fcntl.h>
#include <termios.h>
#include <sys/time.h>
#include <unistd.h>
#include <wiringPi.h>

int g_stop;

#define BAUD_DELAY_45 22000    /* 22ms = 45 baud, 60WPM    */
#define BAUD_DELAY_50 20000    /* 20ms = 50 baud, 66WPM    */
#define BAUD_DELAY_57 18000    /* 18ms = 56.9 baud, 75WPM  */
#define BAUD_DELAY_74 13470    /* 13ms = 74 baud, 100WPM   */
#define COLUMN_MAX 76

#define CHAR_A          0
#define CHAR_Z          25
#define CHAR_NULL       26
#define CHAR_LF         27
#define CHAR_SPACE      28
#define CHAR_CR         29
#define CHAR_SHIFT_UP   30
#define CHAR_SHIFT_DOWN 31
#define CHAR_OPEN       32
#define CHAR_CLOSED     33

#define CHAR_0          15
#define CHAR_1          16
#define CHAR_2          22
#define CHAR_3           4
#define CHAR_4          17
#define CHAR_5          19
#define CHAR_6          24
#define CHAR_7          20
#define CHAR_8           8
#define CHAR_9          14
 
#define CHAR_DASH        0
#define CHAR_QUESTION    1
#define CHAR_COLON       2
#define CHAR_DOLLAR      3
#define CHAR_BELL        6
#define CHAR_APOSTROPHE  9
#define CHAR_LPHAREN    10
#define CHAR_RPHAREN    11
#define CHAR_PERIOD     12
#define CHAR_COMMA      13
#define CHAR_SEMICOLON  21
#define CHAR_SOLIDUS    23
#define CHAR_QUOTE      25

typedef struct rtty_conf
{
    char *filename;
    int bits;
    int bit_delay;
    int wpm;
    int format;
    int shift;
    int column;
    int test_count;
    int sec_sleep;
    int milli_sleep;
    int no_init;
} rtty_conf;

void encode_bit(rtty_conf *ctx, int c);
void encode_to_baudot(rtty_conf *ctx, char c);
void encode_bits(rtty_conf *ctx, char *bits, int len);
int ascii_2_baudot(char c, char *baudot, int *shift);
void print_char(rtty_conf *ctx, char c);
void initialize_tty(rtty_conf *ctx);
void pause_print(rtty_conf *ctx, int count);
void print_line(rtty_conf *ctx, char *line);
void print_file(rtty_conf *ctx, char *name);
int set_raw(int fd, struct termios *old_mode);
void keyboard_io(rtty_conf *ctx);

void interrupt_function(int sig)
{
    if (sig == SIGINT)
    {
        g_stop = 1;
    }
}

void rtty_conf_init(rtty_conf *ctx)
{
    ctx->bits = 8;
    switch (ctx->wpm)
    {
      case 60:
        ctx->bit_delay = BAUD_DELAY_45;
        break;
      case 66:
        ctx->bit_delay = BAUD_DELAY_50;
        break;
      case 75:
        ctx->bit_delay = BAUD_DELAY_57;
        break;
      case 100:
        ctx->bit_delay = BAUD_DELAY_74;
        break;
      default:
        ctx->bit_delay = BAUD_DELAY_45;
        break;
    }

    wiringPiSetup ();
    pinMode (0, OUTPUT);
} 

void
Usage(void) {
    fprintf(stderr, "usage: rtty-alsa [options] number ...\n"
            " Valid options with their default values are:\n"
            "   TTY options:\n"
            "     --input-file IN_FILE\n"
            "     --test-data [N]; default N=1\n"
            "     --keyboard\n"
            "     --wpm 60 | 66 | 100\n"
            "     --char-delay sec\n"
            "     --no-init\n"
            );
            
    exit(1);
}


void
getvalue(int *arg, int *index, int argc,
     char **argv, int min, int max) 
{
    if (*index >= argc-1)
        Usage();

    *arg = atoi(argv[1+*index]);

    if (*arg < min || *arg > max) {
        fprintf(stderr, "Value for %s should be in the range %d..%d\n", 
                argv[*index]+2, min, max);
        exit(1);
    }
    ++*index;
}


void test_generator(rtty_conf *ctx)
{
    int i = 0;
    char line[] =
        "the quick brown fox jumped over the lazy dog's back 1234567890\n"
        "ryryryryryryryryryryryryryryryryryryryryryryryryryryryryryryry";

    initialize_tty(ctx);
    while (ctx->test_count-- > 0) {
        for (i=0; line[i]; i++) {
            print_char(ctx, line[i]);
        }
        print_char(ctx, '\n');
    }
    initialize_tty(ctx);
}


int main(int argc, char **argv)
{
    int i;
    int test_data = 0;
    int keyboard = 0;
    rtty_conf ctx = {0};

    for(i = 1; i < argc; i++) 
    {
        if (argv[i][0] != '-' || argv[i][1] != '-')
            break;

        if (!strcmp(argv[i], "--keyboard")) {
            keyboard = 1;
        }
        else if (!strcmp(argv[i], "--test-data")) {
            test_data = 1;
            ctx.test_count = 1;
            if ((i+1) < argc && atoi(argv[i+1]) > 1) {
                ctx.test_count = atoi(argv[i+1]);
            }
        }
        else if (!strcmp(argv[i], "--wpm")) {
            getvalue(&ctx.wpm, &i, argc, argv,
                 10, 10000);
        }
        else if (!strcmp(argv[i], "--input-file")) {
            i++;
            if (i >= argc)
                Usage();
            ctx.filename = argv[i];
        }
        else if (!strcmp(argv[i], "--no-init")) {
            ctx.no_init = 1;
	}
        else if (!strcmp(argv[i], "--char-delay")) {
            char *b;
            char *e;
            i++;
            if (i >= argc)
                Usage();
            b = argv[i];
            ctx.sec_sleep = strtol(b, &e, 10);
            if (e == b) {
                ctx.sec_sleep = 0;
            }
            else if (e[0] == '.' && e[1] != '\0') {
                b = e+1;
                ctx.milli_sleep = strtol(b, &e, 10);
                if (e == b) {
                    ctx.milli_sleep = 0;
                }
                else {
                    /* 0.1 = 100ms, 0.001 = 1ms */
                    switch (e-b) {
                        case 2:
                            ctx.milli_sleep *= 10;
                            break;
                        case 1:
                            ctx.milli_sleep *= 100;
                    }
                }
            }
            if (ctx.sec_sleep < 0 || ctx.sec_sleep > 10) {
                /* Be nice and give a second betweeen characters */
                ctx.sec_sleep = 1;
                ctx.milli_sleep = 0;
            }
        }
        else {
            Usage();
        }
    }

    rtty_conf_init(&ctx);
    signal(SIGINT, interrupt_function);

    if (test_data) {
        test_generator(&ctx);
    }
    else {
        pause_print(&ctx, 10);
        initialize_tty(&ctx);

        if (keyboard) 
        {
            keyboard_io(&ctx);
        }
        else if (ctx.filename) {
            print_file(&ctx, ctx.filename);
        }
        else {
            while (i < argc) {
                print_line(&ctx, argv[i]);
            }
        }

        initialize_tty(&ctx);
        pause_print(&ctx, 10);
    }
    return 0;
}

void delay_usec(int usec)
{
    struct timeval to = {0};

    to.tv_usec = usec;
    select(0, 0, 0, 0, &to);
}

void keyboard_io(rtty_conf *ctx)
{
    fd_set rmask;
    char c;
    int n;
    struct termios old;
    static int N_count = 0;
    static char prev_char = 0;

    set_raw(0, &old);
    do 
    {
        FD_ZERO(&rmask);
        FD_SET(0, &rmask);
        n = select(1, &rmask, NULL, NULL, NULL);
        if (n > 0) {
            read(0, &c, 1);
            if (c == '\r' || c == '\n' ||
                (ctx->column && ((ctx->column % COLUMN_MAX) == 0)))
            {
                encode_to_baudot(ctx, CHAR_CR);
                encode_to_baudot(ctx, CHAR_LF);
                printf("\r\n");
                ctx->column = 0;
                N_count = 0;
            }
            else 
            {
                print_char(ctx, c);
            }
            fflush(stdout);
            if (c == 'N' && c == prev_char)
            {
                if (ctx->column <= 4)
                {
                    N_count++;
                }
                else
                {
                    N_count = 0;
                }
            }
            prev_char = c;
        }
    } while (N_count < 3);
    tcsetattr(0, TCSANOW, &old);
}

void print_file(rtty_conf *ctx, char *name)
{
    FILE *fp;
    char *line;

    fp = fopen(name, "r");
    if (!fp) {
        perror("fopen");
        return;
    }
    line = (char *) malloc(1024);
    if (!line) {
        perror("malloc");
        fclose(fp);
        return;
    }
    fgets(line, 1022, fp);
    while (!feof(fp)) {
        print_line(ctx, line);
        fgets(line, 1022, fp);
    }
    
    fclose(fp);
    free(line);
}


void print_line(rtty_conf *ctx, char *line)
{
    char *cp;
    if (!line) return;

    for (cp=line; *cp; cp++) {
        if (isspace(*cp) || isalnum(*cp) || ispunct(*cp)) {
            print_char(ctx, *cp);
        }
    }
}


void initialize_tty(rtty_conf *ctx)
{
    if (ctx->no_init) {
        return;
    }
    encode_to_baudot(ctx, CHAR_NULL);
    encode_to_baudot(ctx, CHAR_NULL);
    encode_to_baudot(ctx, CHAR_SHIFT_DOWN);
    encode_to_baudot(ctx, CHAR_CR);
    encode_to_baudot(ctx, CHAR_LF);
}


void pause_print(rtty_conf *ctx, int count)
{
    int i;
    for (i=0; i<count; i++) {
        encode_to_baudot(ctx, CHAR_CLOSED);
    }
}


void print_char(rtty_conf *ctx, char c)
{
    char baudot[16];
    char *bp;
    int cnt;

    if (g_stop)
    {
        initialize_tty(ctx);
        encode_to_baudot(ctx, CHAR_CLOSED);
        exit(0);
    }

    cnt = ascii_2_baudot(c, baudot, &ctx->shift);
    bp = baudot;
    while (cnt > 0) {
        encode_to_baudot(ctx, *bp);
        bp++;
        cnt--;
    }

    if (isspace(c) || isalnum(c) || ispunct(c))
    {
        ctx->column++;
        if (c == '\n' || c == '\r')
        {
            ctx->column = 0;
            printf("\r\n");
        }
        else
        {
            putchar((char) toupper((int) c));
        }
        fflush(stdout);
    }

    if (ctx->column >= COLUMN_MAX) {
        encode_to_baudot(ctx, CHAR_CR);
        encode_to_baudot(ctx, CHAR_LF);
        encode_to_baudot(ctx, CHAR_CR);
        printf("\r\n");
        ctx->column = 0;
    }
    if (ctx->sec_sleep > 0 || ctx->milli_sleep > 0) {
        sleep(ctx->sec_sleep);
        usleep(ctx->milli_sleep * 1000);
    }
}


int
ascii_2_baudot(char c, char *baudot, int *shift)
{
    int i;
    char *p = baudot;

    /* Convert punction characters first */
    switch (c) {
      case '-':
        *p++ = CHAR_DASH;
        break;
      case '?':
        *p++ = CHAR_QUESTION;
        break;
      case ':':
        *p++ = CHAR_COLON;
        break;
      case '$':
        *p++ = CHAR_DOLLAR;
        break;
      case 7:
        *p++ = CHAR_BELL;
        break;
      case '\'':
      case '`':
        *p++ = CHAR_APOSTROPHE;
        break;
      case '(':
        *p++ = CHAR_LPHAREN;
        break;
      case ')':
        *p++ = CHAR_RPHAREN;
        break;
      case '.':
        *p++ = CHAR_PERIOD;
        break;
      case ',':
        *p++ = CHAR_COMMA;
        break;
      case ';':
        *p++ = CHAR_SEMICOLON;
        break;
      case '/':
        *p++ = CHAR_SOLIDUS;
        break;
      case '"':
        *p++ = CHAR_QUOTE;
        break;
    }

    /* p has advanced one byte if we have found a character to convert */

    if (p != baudot) {
        /*
         * Prefix the converted character with shift up if not 
         * already shifted
         */
        if (!*shift) {
            p[0] = p[-1];
            p[-1] = CHAR_SHIFT_UP;
            p++;
            *shift = 1;
        }
    }
    else {

        /* The current character has not been mapped */

        i = (int) c;
        if (c == ' ') {
            *p++ = CHAR_SPACE;
        }
        else if (c == '\n') {
            *p++ = CHAR_CR;
            *p++ = CHAR_LF;
        }
        else if (isdigit(c)) {
            if (!*shift) {
                *shift = 1;
                *p++ = CHAR_SHIFT_UP;
            }
            i = (int) (c - '0');
            switch(i) {
              case 0:
                *p++ = CHAR_0;
                break;
              case 1:
                *p++ = CHAR_1;
                break;
              case 2:
                *p++ = CHAR_2;
                break;
              case 3:
                *p++ = CHAR_3;
                break;
              case 4:
                *p++ = CHAR_4;
                break;
              case 5:
                *p++ = CHAR_5;
                break;
              case 6:
                *p++ = CHAR_6;
                break;
              case 7:
                *p++ = CHAR_7;
                break;
              case 8:
                *p++ = CHAR_8;
                break;
              case 9:
                *p++ = CHAR_9;
                break;
            }
        }
        else if (isalpha(c)) {
            if (*shift) {
                *shift = 0;
                *p++ = CHAR_SHIFT_DOWN;
            }
            *p++ = (char) (toupper(c) - 'A');
        }
        else if (i>=CHAR_NULL && i<=CHAR_CLOSED) {
           *p++ = i;
        }
        else {
            /*
             * There is no reasonable mapping for this character
             */
            *p++ = CHAR_NULL;
        }
    }
    return p - baudot;
}


void
encode_to_baudot(rtty_conf *ctx, char c)
{
    static char baudot_bits[][8] = {
        {0, 1, 1, 0, 0, 0, 1, 1}, /* A */
        {0, 1, 0, 0, 1, 1, 1, 1}, /* B */
        {0, 0, 1, 1, 1, 0, 1, 1}, /* C */
        {0, 1, 0, 0, 1, 0, 1, 1}, /* D */
        {0, 1, 0, 0, 0, 0, 1, 1}, /* E / 3 */
        {0, 1, 0, 1, 1, 0, 1, 1}, /* F */
        {0, 0, 1, 0, 1, 1, 1, 1}, /* G */
        {0, 0, 0, 1, 0, 1, 1, 1}, /* H */
        {0, 0, 1, 1, 0, 0, 1, 1}, /* I  / 8 */
        {0, 1, 1, 0, 1, 0, 1, 1}, /* J */
        {0, 1, 1, 1, 1, 0, 1, 1}, /* K */
        {0, 0, 1, 0, 0, 1, 1, 1}, /* L */
        {0, 0, 0, 1, 1, 1, 1, 1}, /* M / . */
        {0, 0, 0, 1, 1, 0, 1, 1}, /* N */
        {0, 0, 0, 0, 1, 1, 1, 1}, /* O / 9 */
        {0, 0, 1, 1, 0, 1, 1, 1}, /* P / 0 */
        {0, 1, 1, 1, 0, 1, 1, 1}, /* Q / 1 */
        {0, 0, 1, 0, 1, 0, 1, 1}, /* R / 4 */
        {0, 1, 0, 1, 0, 0, 1, 1}, /* S */
        {0, 0, 0, 0, 0, 1, 1, 1}, /* T / 5 */
        {0, 1, 1, 1, 0, 0, 1, 1}, /* U / 7 */
        {0, 0, 1, 1, 1, 1, 1, 1}, /* V */
        {0, 1, 1, 0, 0, 1, 1, 1}, /* W / 2 */
        {0, 1, 0, 1, 1, 1, 1, 1}, /* X / / */
        {0, 1, 0, 1, 0, 1, 1, 1}, /* Y / 6 */
        {0, 1, 0, 0, 0, 1, 1, 1}, /* Z */
        {0, 0, 0, 0, 0, 0, 1, 1}, /* NULL */
        {0, 0, 1, 0, 0, 0, 1, 1}, /* LF */
        {0, 0, 0, 1, 0, 0, 1, 1}, /* SPACE */
        {0, 0, 0, 0, 1, 0, 1, 1}, /* CR */
        {0, 1, 1, 0, 1, 1, 1, 1}, /* SHIFT_UP */
        {0, 1, 1, 1, 1, 1, 1, 1}, /* SHIFT_DOWN */
        {0, 0, 0, 0, 0, 0, 0, 0}, /* Open */
        {1, 1, 1, 1, 1, 1, 1, 1}  /* closed */
    };
    int i = (int) c;

    if (i>=0 && i<(sizeof(baudot_bits)/8)) {
       encode_bits(ctx, baudot_bits[i], 8);
    }
}

void encode_bits(rtty_conf *ctx, char *bits, int len)
{
    int i;

    for (i=0; i < len; i++) {
        encode_bit(ctx, bits[i]);
    }
}

void
encode_bit(rtty_conf *ctx, int c) 
{
    switch(c) {
    case 0:
        digitalWrite(0, HIGH); 
        delay_usec(ctx->bit_delay);
        break;
    case 1:
        digitalWrite(0, LOW); 
        delay_usec(ctx->bit_delay);
        break;
    }
}

/*
 *  Put terminal into raw mode.
 *  The "definitive" set raw code?  I don't know about that
 *  but it works.
 *
 *  Until proven otherwise this is the way to put a pty into raw mode.
 *  This has been tested by streaming binary data into and from
 *  a pty after being set into raw mode with this code, and it
 *  was completely unmodified by the pty.  Earlier versions of
 *  set_raw() could not claim this fact.
 */
int set_raw(int fd, struct termios *old_mode)
{
    struct termios mode;

    if (!isatty(fd)) {
        return (-1);
    }

    tcgetattr(fd, &mode);
    if (old_mode) {
        tcgetattr(fd, old_mode);
    }

    mode.c_iflag     = 0;
    mode.c_oflag    &= ~OPOST;
#ifdef NO_XCASE
    mode.c_lflag    &= ~(IEXTEN | ISIG | ICANON | ECHO);
#else
    mode.c_lflag    &= ~(IEXTEN | ISIG | ICANON | ECHO | XCASE);
#endif
    mode.c_cflag    &= ~(CSIZE | PARENB);
    mode.c_cflag    |= CS8;
    mode.c_cc[VMIN]  = 1;
    mode.c_cc[VTIME] = 1;

    tcsetattr(fd, TCSANOW, &mode);
    return 0;
}
